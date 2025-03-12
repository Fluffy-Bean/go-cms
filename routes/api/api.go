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
	"github.com/google/uuid"
)

func RegisterAPIRoutes(mux *http.ServeMux, h *handler.Handler) {
	mux.HandleFunc("/api/v1/page:create", routePageCreate(h))
	mux.HandleFunc("/api/v1/page:delete", routePageDelete(h))

	mux.HandleFunc("/api/v1/blocks:available", routeBlocksAvailable(h))
}

func routePageCreate(h *handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		if strings.TrimSpace(r.Form.Get("core.page_url")) == "" {
			fmt.Println("missing url")
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		ID := r.URL.Query().Get("id")

		var page handler.Page
		if ID != "" {
			page, err = h.GetPage(ID)
			if err != nil {
				fmt.Println(err)
				http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

				return
			}
		} else {
			page, err = h.NewPage()
			if err != nil {
				fmt.Println(err)
				http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

				return
			}
		}

		page.Path = strings.TrimSpace(r.Form.Get("core.page_url"))
		page.Meta.Title = strings.TrimSpace(r.Form.Get("core.page_title"))
		page.Meta.Description = strings.TrimSpace(r.Form.Get("core.page_description"))

		blockIndex := 0
		blockFormData := map[string]map[string]string{}
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

			if _, exists := blockFormData[fieldEntries[1]]; !exists {
				blockFormData[fieldEntries[1]] = map[string]string{
					// Should somehow avoid this...
					"Index": fmt.Sprintf("%d", blockIndex),
				}

				blockIndex += 1
			}

			blockFormData[fieldEntries[1]][fieldEntries[2]] = value[0]
		}

		pageHTML := make([]string, blockIndex)
		pageBlocks := make([]blocks.Handle, blockIndex)
		for blockID, formData := range blockFormData {
			_index, ok := formData["Index"]
			if !ok {
				fmt.Println(err)
				http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

				return
			}
			index, err := strconv.Atoi(_index)
			if err != nil {
				fmt.Println(err)
				http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

				return
			}

			block, err := h.Blocks.GetBlock(blockID)
			if err != nil {
				fmt.Println(err)
				http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

				return
			}

			block, err = h.Blocks.ParseFormIntoBlock(formData, block)
			if err != nil {
				fmt.Println(err)
				http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

				return
			}

			html := h.Blocks.Render(block)

			pageHTML = slices.Insert(pageHTML, index, html)
			pageBlocks = slices.Insert(pageBlocks, index, block)
		}

		if page.TemplateID != "" {
			err = os.Remove(h.DataPath + "/pages/" + page.TemplateID)
			if err != nil {
				fmt.Println(err)
				http.Redirect(w, r, "/cms/pages?status=failure", http.StatusFound)

				return
			}
		}

		newTemplateID := uuid.New().String() + ".html"

		templateFile, err := os.Create(h.DataPath + "/pages/" + newTemplateID)
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
			"Title":       page.Meta.Title,
			"Description": page.Meta.Description,
			"Blocks":      pageHTML,
		})
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		page.TemplateID = newTemplateID
		page.Blocks = pageBlocks

		err = h.UpdatePage(page)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		http.Redirect(w, r, "/cms/pages?status=success", http.StatusFound)
	}
}

func routePageDelete(h *handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageURL := r.URL.Query().Get("page")

		page, err := h.GetPage(pageURL)
		if err != nil {
			http.NotFound(w, r)

			return
		}

		err = os.Remove(h.DataPath + "/pages/" + page.TemplateID)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/pages?status=failure", http.StatusFound)

			return
		}

		h.DeletePage(page)

		http.Redirect(w, r, "/cms/pages?status=success", http.StatusFound)
	}
}

// Its like 0020 and I cant be asked to get this working currently...
func routeBlocksAvailable(h *handler.Handler) http.HandlerFunc {
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
