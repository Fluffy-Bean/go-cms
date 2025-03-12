package api

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/Fluffy-Bean/cms/internal/blocks"
	"github.com/Fluffy-Bean/cms/internal/handler"
	"github.com/Fluffy-Bean/cms/internal/router"
	"github.com/google/uuid"
)

func RegisterAPIRoutes(mux *http.ServeMux, h handler.Handler) {
	mux.HandleFunc("/api/v1/page:create", routePageCreate(h))
	mux.HandleFunc("/api/v1/page:delete", routePageDelete(h))

	mux.HandleFunc("/api/v1/blocks:available", routeBlocksAvailable(h))
}

func routePageCreate(h handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		pageID := r.URL.Query().Get("id")

		var route router.Route
		route, err = h.Router.GetRoute(pageID)
		if err != nil {
			route, err = h.Router.NewRoute()
			if err != nil {
				fmt.Println(err)
				http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

				return
			}
		}

		route.Path = strings.TrimSpace(r.Form.Get("core.page_url"))
		route.Meta.Title = strings.TrimSpace(r.Form.Get("core.page_title"))
		route.Meta.Description = strings.TrimSpace(r.Form.Get("core.page_description"))

		if route.Path == "" {
			fmt.Println("missing url")
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		index := 0
		formData := map[string]map[string]string{}
		// ID:
		//   Index: "0"
		//   Field: "Value"
		//   Field: "Value"
		//
		// ID:
		//   Index: "1"
		//   Field: "Value"
		//   Field: "Value"
		//
		for field, value := range r.Form {
			if !strings.HasPrefix(field, "block.") {
				continue
			}

			fieldEntries := strings.Split(field, ".")
			if len(fieldEntries) != 3 {
				continue
			}

			// 1 - ID
			// 2 - Struct Field

			if _, exists := formData[fieldEntries[1]]; !exists {
				formData[fieldEntries[1]] = map[string]string{
					// Should somehow avoid this...
					"Index": fmt.Sprintf("%d", index),
				}

				index += 1
			}

			formData[fieldEntries[1]][fieldEntries[2]] = value[0]
		}

		pageHTML := make([]string, len(formData))
		pageBlocks := make([]blocks.Handle, len(formData))
		for id, data := range formData {
			_index, ok := data["Index"]
			if !ok {
				continue
			}
			index, err := strconv.Atoi(_index)
			if err != nil {
				fmt.Println(err)

				continue
			}

			block, err := h.Blocks.GetBlock(id)
			if err != nil {
				fmt.Println(err)
				http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

				return
			}

			block, err = h.Blocks.ParseFormIntoBlock(data, block)
			if err != nil {
				fmt.Println(err)

				continue
			}

			html := h.Blocks.Render(block)

			pageHTML = slices.Insert(pageHTML, index, html)
			pageBlocks = slices.Insert(pageBlocks, index, block)
		}

		if route.TemplateID != "" {
			err = os.Remove(h.DataPath + "/routes/" + route.TemplateID)
			if err != nil {
				fmt.Println(err)
				http.Redirect(w, r, "/cms/pages?status=failure", http.StatusFound)

				return
			}
		}

		newTemplateID := uuid.New().String() + ".html"

		templateFile, err := os.Create(h.DataPath + "/routes/" + newTemplateID)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		templ, err := template.ParseFiles(h.TemplatesPath + "/shell.html")
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		err = templ.Execute(templateFile, map[string]any{
			"Title":       route.Meta.Title,
			"Description": route.Meta.Description,
			"Blocks":      pageHTML,
		})
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		route.TemplateID = newTemplateID
		route.Blocks = pageBlocks

		err = h.Router.UpdateRoute(route)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		http.Redirect(w, r, "/cms/pages?status=success", http.StatusFound)
	}
}

func routePageDelete(h handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageURL := r.URL.Query().Get("page")

		page, err := h.Router.GetRoute(pageURL)
		if err != nil {
			http.NotFound(w, r)

			return
		}

		err = os.Remove(h.DataPath + "/routes/" + page.TemplateID)
		if err != nil {
			fmt.Println(err)

			http.Redirect(w, r, "/cms/pages?status=failure", http.StatusFound)

			return
		}

		h.Router.DeleteRoute(page)

		http.Redirect(w, r, "/cms/pages?status=success", http.StatusFound)
	}
}

// Its like 0020 and I cant be asked to get this working currently...
func routeBlocksAvailable(h handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//var form []blocks.FormData
		//for _, block := range h.Blocks.GetRegisteredBlocksIDs() {
		//    formData, err := h.Blocks.GetFormData(block)
		//    if err != nil {
		//        fmt.Println(err)
		//
		//        continue
		//    }
		//
		//    form = append(form, formData)
		//}
		//
		//bytes, err := json.Marshal(form)
		//if err != nil {
		//    fmt.Println(err)
		//    http.Error(w, err.Error(), http.StatusInternalServerError)
		//
		//    return
		//}
		//
		//w.Header().Set("Content-Type", "application/json; charset=utf-8")
		//
		//w.Write(bytes)
	}
}
