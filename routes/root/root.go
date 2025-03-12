package root

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Fluffy-Bean/cms/internal/handler"
)

func RegisterRootRoutes(mux *http.ServeMux, h *handler.Handler) {
	mux.HandleFunc("/", routeRootHandler(h))
}

func routeRootHandler(h *handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := h.GetPage(r.URL.Path)
		if err != nil {
			http.NotFound(w, r)

			return
		}

		templ, err := template.ParseFiles(
			h.DataPath+"/pages/"+page.TemplateID,
			h.TemplatesPath+"/generated.html",
		)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = templ.Execute(w, map[string]any{
			"Page": page,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}
}
