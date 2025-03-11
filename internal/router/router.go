package router

import (
	"errors"

	"github.com/Fluffy-Bean/cms/internal/blocks"
)

type Route struct {
	Meta struct {
		Title       string
		Description string
	}
	TemplateID string
	Blocks     []struct {
		ID    string
		Block blocks.Block
	}
}

type Router struct {
	Routes map[string]Route
}

func New() Router {
	return Router{
		Routes: map[string]Route{},
	}
}

func (r *Router) RegisterRoute(path string, route Route) error {
	if path == "" {
		return errors.New("empty path")
	}
	if _, ok := r.Routes[path]; ok {
		return errors.New("route already exists")
	}

	r.Routes[path] = route

	return nil
}

func (r *Router) FindRoute(path string) (Route, error) {
	route, ok := r.Routes[path]
	if !ok {
		return Route{}, errors.New("route not found")
	}

	return route, nil
}

func (r *Router) RemoveRoute(path string) {
	delete(r.Routes, path)
}
