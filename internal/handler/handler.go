package handler

import (
	"github.com/Fluffy-Bean/cms/internal/blocks"
	"github.com/Fluffy-Bean/cms/internal/router"
)

type Handler struct {
	TemplatesPath string
	DataPath      string
	Router        router.Router
	Blocks        blocks.Blocks
}
