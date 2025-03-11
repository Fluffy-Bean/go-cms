package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Fluffy-Bean/cms/internal/blocks"
	"github.com/Fluffy-Bean/cms/internal/handler"
	"github.com/Fluffy-Bean/cms/internal/router"
	"github.com/Fluffy-Bean/cms/routes/api"
	"github.com/Fluffy-Bean/cms/routes/cms"
	"github.com/Fluffy-Bean/cms/routes/root"
)

type Config struct {
	Host string
}

func Execute(conf Config) {
	h := handler.Handler{
		TemplatesPath: "./templates",
		DataPath:      "./_data",
		Router:        router.New(),
		Blocks:        blocks.New(),
	}

	_ = h.Blocks.RegisterBlock("core:text", blocks.TextBlock{})
	_ = h.Blocks.RegisterBlock("core:blogPost", blocks.BlogPostBlock{})
	_ = h.Blocks.RegisterBlock("core:code", blocks.CodeBlock{})
	_ = h.Blocks.RegisterBlock("core:image", blocks.ImageBlock{})

	mux := http.NewServeMux()

	cms.RegisterCMSRoutes(mux, h)
	api.RegisterAPIRoutes(mux, h)
	root.RegisterRootRoutes(mux, h)

	fmt.Println("Serving on: ", conf.Host)
	log.Fatal(http.ListenAndServe(conf.Host, mux))
}
