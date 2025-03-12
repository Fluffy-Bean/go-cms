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

func RegisterCMSRoutes(mux *http.ServeMux, h *handler.Handler) {
	mux.HandleFunc("/static/", routeStatic(h))

	mux.HandleFunc("/cms", routeRoot(h))
	mux.HandleFunc("/cms/editor", routeEditor(h))
	mux.HandleFunc("/cms/pages", routePages(h))
	mux.HandleFunc("/cms/profile", routeProfile(h))
	mux.HandleFunc("/cms/files", routeFiles(h))
}

func routeStatic(h *handler.Handler) http.HandlerFunc {
	cssStyles, _ := os.ReadFile("./static/css/styles.css")
	cssBlocks, _ := os.ReadFile("./static/css/blocks.css")
	jsDOM, _ := os.ReadFile("./static/js/dom.js")

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/static/css/styles.css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			w.Write(cssStyles)
		case "/static/css/blocks.css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			w.Write(cssBlocks)
		case "/static/js/dom.js":
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			w.Write(jsDOM)
		default:
			http.NotFound(w, r)
		}
	}
}

func routeRoot(h *handler.Handler) http.HandlerFunc {
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

func routeEditor(h *handler.Handler) http.HandlerFunc {
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

		pageID := ""
		pageUrl := ""
		pageTitle := ""
		pageDescription := ""
		pageBlocks := make([]blocks.FormData, 0)

		if page != "" {
			route, err := h.GetPage(page)
			if err != nil {
				fmt.Println(err)

				http.NotFound(w, r)

				return
			}

			pageID = route.ID
			pageUrl = route.Path
			pageTitle = route.Meta.Title
			pageDescription = route.Meta.Description

			for _, block := range route.Blocks {
				formData, err := h.Blocks.GetFormData(block)
				if err != nil {
					fmt.Println(err)

					continue
				}
				pageBlocks = append(pageBlocks, formData)
			}
		}

		if slots != "" {
			for _, id := range strings.Split(slots, ",") {
				block, err := h.Blocks.NewBlock(id)
				if err != nil {
					fmt.Println(err)

					continue
				}

				formData, err := h.Blocks.GetFormData(block)
				if err != nil {
					fmt.Println(err)

					continue
				}

				pageBlocks = append(pageBlocks, formData)
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = templ.Execute(w, map[string]any{
			"PageID":          pageID,
			"PageURL":         pageUrl,
			"PageTitle":       pageTitle,
			"PageDescription": pageDescription,
			"Message":         message,
			"Blocks":          pageBlocks,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

func routePages(h *handler.Handler) http.HandlerFunc {
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
			"Pages":   h.Pages,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routeProfile(h *handler.Handler) http.HandlerFunc {
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

func routeFiles(h *handler.Handler) http.HandlerFunc {
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
