package cms

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/Fluffy-Bean/cms/app"
	"github.com/Fluffy-Bean/cms/internal/blocks"
)

func RegisterCMSRoutes(mux *http.ServeMux, handler app.App) {
	mux.HandleFunc("/static/", routeStatic(handler))

	mux.HandleFunc("/cms", routeRoot(handler))
	mux.HandleFunc("/cms/editor", routeEditor(handler))
	mux.HandleFunc("/cms/pages", routePages(handler))
	mux.HandleFunc("/cms/profile", routeProfile(handler))
	mux.HandleFunc("/cms/files", routeFiles(handler))
}

func routeStatic(handler app.App) http.HandlerFunc {
	css, _ := os.ReadFile("./static/css/styles.css")
	blocks, _ := os.ReadFile("./static/css/blocks.css")
	js, _ := os.ReadFile("./static/js/dom.js")

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/static/css/styles.css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			w.Write(css)

			return
		case "/static/css/blocks.css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			w.Write(blocks)

			return
		case "/static/js/dom.js":
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			w.Write(js)

			return
		default:
			http.NotFound(w, r)
		}
	}
}

func routeRoot(handler app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles(
			handler.TemplatesPath+"/cms/root.html",
			handler.TemplatesPath+"/cms.html",
		)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = templ.Execute(w, nil)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routeEditor(handler app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		page := r.URL.Query().Get("page")
		slots := r.URL.Query().Get("slots")

		var message string
		switch status {
		case "success":
			message = "Success"
		case "failure":
			message = "Failure"
		}

		templ, err := template.ParseFiles(
			handler.TemplatesPath+"/cms/editor.html",
			handler.TemplatesPath+"/cms.html",
		)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		pageUrl := ""
		pageTitle := ""
		pageDescription := ""
		pageNew := true

		var form []blocks.FormData
		if page != "" {
			pageData, err := handler.Router.FindRoute(page)
			if err != nil {
				fmt.Println(err)

				http.NotFound(w, r)

				return
			}

			pageUrl = page
			pageTitle = pageData.Meta.Title
			pageDescription = pageData.Meta.Description
			pageNew = false

			for index, block := range pageData.Blocks {
				formData, err := handler.Blocks.GetFormDataByType(block.Block)
				if err != nil {
					fmt.Println(err)

					continue
				}
				formData.Index = index
				formData.ID = block.ID

				form = append(form, formData)
			}
		}

		indexOffset := len(form)
		if slots != "" {
			for index, block := range strings.Split(slots, ",") {
				formData, err := handler.Blocks.GetFormDataByID(block)
				if err != nil {
					fmt.Println(err)

					continue
				}
				formData.Index = index + indexOffset

				form = append(form, formData)
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = templ.Execute(w, map[string]any{
			"NewPage":     pageNew,
			"URL":         pageUrl,
			"Title":       pageTitle,
			"Description": pageDescription,
			"Message":     message,
			"Blocks":      form,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

func routePages(handler app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")

		var message string
		switch status {
		case "success":
			message = "Success"
		case "failure":
			message = "Failure"
		}

		templ, err := template.ParseFiles(
			handler.TemplatesPath+"/cms/pages.html",
			handler.TemplatesPath+"/cms.html",
		)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = templ.Execute(w, map[string]any{
			"Message": message,
			"Routes":  handler.Router.Routes,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routeProfile(handler app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles(
			handler.TemplatesPath+"/cms/profile.html",
			handler.TemplatesPath+"/cms.html",
		)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = templ.Execute(w, nil)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routeFiles(handler app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles(
			handler.TemplatesPath+"/cms/files.html",
			handler.TemplatesPath+"/cms.html",
		)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = templ.Execute(w, nil)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
