package root

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Fluffy-Bean/cms/app"
)

func RegisterRootRoutes(mux *http.ServeMux, handler app.App) {
	mux.HandleFunc("/", routeRootHandler(handler))
}

func routeRootHandler(handler app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route, err := handler.Router.FindRoute(r.URL.Path)
		if err != nil {
			http.NotFound(w, r)

			return
		}

		templ, err := template.ParseFiles(
			handler.DataPath+"/routes/"+route.TemplateID,
			handler.TemplatesPath+"/generated.html",
		)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = templ.Execute(w, map[string]any{
			"URL":         r.URL.Path,
			"Title":       route.Meta.Title,
			"Description": route.Meta.Description,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}
}
