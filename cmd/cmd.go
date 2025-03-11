package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Fluffy-Bean/cms/app"
	"github.com/Fluffy-Bean/cms/internal/blocks"
	"github.com/Fluffy-Bean/cms/internal/router"
	"github.com/Fluffy-Bean/cms/routes/api"
	"github.com/Fluffy-Bean/cms/routes/cms"
	"github.com/Fluffy-Bean/cms/routes/root"
)

type Config struct {
	Host string
}

func Execute(conf Config) {
	handler := app.App{
		TemplatesPath: "./templates",
		DataPath:      "./_data",
		Router:        router.New(),
		Blocks:        blocks.New(),
	}

	handler.Blocks.RegisterBlock("core:text", blocks.TextBlock{})
	handler.Blocks.RegisterBlock("core:blogPost", blocks.BlogPostBlock{})
	handler.Blocks.RegisterBlock("core:code", blocks.CodeBlock{})
	handler.Blocks.RegisterBlock("core:image", blocks.ImageBlock{})

	mux := http.NewServeMux()

	cms.RegisterCMSRoutes(mux, handler)
	api.RegisterAPIRoutes(mux, handler)
	root.RegisterRootRoutes(mux, handler)

	fmt.Println("Serving on: ", conf.Host)
	log.Fatal(http.ListenAndServe(conf.Host, mux))
}
