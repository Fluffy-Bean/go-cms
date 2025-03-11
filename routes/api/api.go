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

	"github.com/Fluffy-Bean/cms/app"
	"github.com/Fluffy-Bean/cms/internal/blocks"
	"github.com/Fluffy-Bean/cms/internal/router"
	"github.com/google/uuid"
)

func RegisterAPIRoutes(mux *http.ServeMux, handler app.App) {
	mux.HandleFunc("/api/v1/page:create", routePageCreate(handler))
	mux.HandleFunc("/api/v1/page:delete", routePageDelete(handler))

	mux.HandleFunc("/api/v1/blocks:available", routeBlocksAvailable(handler))
}

func routePageCreate(handler app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		pageEditing := r.URL.Query().Get("editing")

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

		if pageEditing == "yeah" {
			page, err := handler.Router.FindRoute(pageURL)
			if err != nil {
				fmt.Println(err)

				http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

				return
			}

			err = os.Remove(handler.DataPath + "/routes/" + page.TemplateID)
			if err != nil {
				fmt.Println(err)

				http.Redirect(w, r, "/cms/pages?status=failure", http.StatusFound)

				return
			}

			handler.Router.RemoveRoute(pageURL)
		}

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

			block, err := handler.Blocks.ParseForm(id, data)
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

		fileName := uuid.New().String() + ".html"

		file, err := os.Create(handler.DataPath + "/routes/" + fileName)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		templ, err := template.ParseFiles(handler.TemplatesPath + "/shell.html")
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		err = templ.Execute(file, map[string]any{
			"Title":       pageTitle,
			"Description": pageDescription,
			"Blocks":      pageHTML,
		})
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		newRoute := router.Route{
			Meta: struct {
				Title       string
				Description string
			}{
				Title:       pageTitle,
				Description: pageDescription,
			},
			TemplateID: fileName,
			Blocks:     pageBlocks,
		}

		err = handler.Router.RegisterRoute(pageURL, newRoute)
		if err != nil {
			http.Redirect(w, r, "/cms/editor?status=failure", http.StatusFound)

			return
		}

		if pageEditing == "yeah" {
			http.Redirect(w, r, "/cms/editor?page="+pageURL+"&status=success", http.StatusFound)
		} else {
			http.Redirect(w, r, "/cms/editor?status=success", http.StatusFound)
		}
	}
}

func routePageDelete(handler app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageURL := r.URL.Query().Get("page")

		page, err := handler.Router.FindRoute(pageURL)
		if err != nil {
			http.NotFound(w, r)

			return
		}

		err = os.Remove(handler.DataPath + "/routes/" + page.TemplateID)
		if err != nil {
			fmt.Println(err)

			http.Redirect(w, r, "/cms/pages?status=failure", http.StatusFound)

			return
		}

		handler.Router.RemoveRoute(pageURL)

		http.Redirect(w, r, "/cms/pages?status=success", http.StatusFound)
	}
}

func routeBlocksAvailable(handler app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form []blocks.FormData
		for _, block := range handler.Blocks.GetRegisteredBlocks() {
			formData, err := handler.Blocks.GetFormDataByID(block)
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
