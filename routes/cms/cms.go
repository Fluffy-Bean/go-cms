package cms

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/Fluffy-Bean/cms/internal/blocks"
	"github.com/Fluffy-Bean/cms/internal/handler"
)

func RegisterCMSRoutes(mux *http.ServeMux, h handler.Handler) {
	mux.HandleFunc("/static/", routeStatic(h))

	mux.HandleFunc("/cms", routeRoot(h))
	mux.HandleFunc("/cms/editor", routeEditor(h))
	mux.HandleFunc("/cms/pages", routePages(h))
	mux.HandleFunc("/cms/profile", routeProfile(h))
	mux.HandleFunc("/cms/files", routeFiles(h))
}

func routeStatic(h handler.Handler) http.HandlerFunc {
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

func routeRoot(h handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles(
			h.TemplatesPath+"/cms/root.html",
			h.TemplatesPath+"/cms.html",
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

func routeEditor(h handler.Handler) http.HandlerFunc {
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
			h.TemplatesPath+"/cms/editor.html",
			h.TemplatesPath+"/cms.html",
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
			pageData, err := h.Router.FindRoute(page)
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
				formData, err := h.Blocks.GetFormDataByType(block.Block)
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
				formData, err := h.Blocks.GetFormDataByID(block)
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
			"Store":       form,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

func routePages(h handler.Handler) http.HandlerFunc {
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
			h.TemplatesPath+"/cms/pages.html",
			h.TemplatesPath+"/cms.html",
		)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = templ.Execute(w, map[string]any{
			"Message": message,
			"Routes":  h.Router.Store,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routeProfile(h handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles(
			h.TemplatesPath+"/cms/profile.html",
			h.TemplatesPath+"/cms.html",
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

func routeFiles(h handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles(
			h.TemplatesPath+"/cms/files.html",
			h.TemplatesPath+"/cms.html",
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
