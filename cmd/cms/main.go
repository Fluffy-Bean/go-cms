package main

import (
	"log"
	"net/http"

	"github.com/Fluffy-Bean/cms/internal/blocks"
	"github.com/Fluffy-Bean/cms/internal/handler"
	"github.com/Fluffy-Bean/cms/routes/api"
	"github.com/Fluffy-Bean/cms/routes/cms"
	"github.com/Fluffy-Bean/cms/routes/root"
)

func main() {
	h := handler.Handler{
		TemplatesPath: "./templates",
		DataPath:      "./_data",
		Blocks:        blocks.New(),
		Pages:         map[string]handler.Page{},
	}

	_ = h.Blocks.RegisterBlock("core:text", blocks.TextBlock{})
	_ = h.Blocks.RegisterBlock("core:blogPost", blocks.BlogPostBlock{})
	_ = h.Blocks.RegisterBlock("core:code", blocks.CodeBlock{})
	_ = h.Blocks.RegisterBlock("core:image", blocks.ImageBlock{})

	mux := http.NewServeMux()

	cms.RegisterCMSRoutes(mux, &h)
	api.RegisterAPIRoutes(mux, &h)
	root.RegisterRootRoutes(mux, &h)

	log.Fatal(http.ListenAndServe(":7070", mux))
}
