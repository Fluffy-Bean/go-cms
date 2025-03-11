package api

import (
	"encoding/json"
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
		pageURL := ""
		pageTitle := ""
		pageDescription := ""
		pageFormData := map[string]map[string]string{}

		for field, data := range r.Form {
			switch field {
			case "core.page_url":
				pageURL = strings.TrimSpace(data[0])
			case "core.page_title":
				pageTitle = strings.TrimSpace(data[0])
			case "core.page_description":
				pageDescription = strings.TrimSpace(data[0])
			default:
				if !strings.HasPrefix(field, "block.") {
					continue
				}

				// 0 - Always "block"
				// 1 - ID
				// 2 - Type's Field
				fieldEntries := strings.Split(field, ".")

				if len(fieldEntries) != 3 {
					continue
				}

				if _, exists := pageFormData[fieldEntries[1]]; !exists {
					pageFormData[fieldEntries[1]] = map[string]string{}
				}

				pageFormData[fieldEntries[1]][fieldEntries[2]] = data[0]
			}
		}

		if pageURL == "" {
			fmt.Println("missing url")
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

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

		route.Path = pageURL
		route.Meta.Title = pageTitle
		route.Meta.Description = pageDescription

		pageHTML := make([]string, len(pageFormData))
		pageBlocks := make([]struct {
			ID    string
			Block blocks.Block
		}, len(pageFormData))
		for i, data := range pageFormData {
			index, err := strconv.Atoi(i)
			if err != nil {
				fmt.Println(err)

				continue
			}
			id, ok := data["ID"]
			if !ok {
				continue
			}

			block, err := h.Blocks.ParseForm(id, data)
			if err != nil {
				fmt.Println(err)

				continue
			}

			html := block.Render()

			pageHTML = slices.Insert(pageHTML, index, html)
			pageBlocks = slices.Insert(pageBlocks, index, struct {
				ID    string
				Block blocks.Block
			}{
				ID:    id,
				Block: block,
			})
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
			"Title":       pageTitle,
			"Description": pageDescription,
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

func routeBlocksAvailable(h handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form []blocks.FormData
		for _, block := range h.Blocks.GetRegisteredBlocks() {
			formData, err := h.Blocks.GetFormDataByID(block)
			if err != nil {
				fmt.Println(err)

				continue
			}

			form = append(form, formData)
		}

		bytes, err := json.Marshal(form)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		w.Write(bytes)
	}
}
