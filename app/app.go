package app

import (
	"github.com/Fluffy-Bean/cms/internal/blocks"
	"github.com/Fluffy-Bean/cms/internal/router"
)

type App struct {
	TemplatesPath string
	DataPath      string
	Router        router.Router
	Blocks        blocks.Blocks
}
