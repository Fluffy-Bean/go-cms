package router

import (
	"errors"

	"github.com/Fluffy-Bean/cms/internal/blocks"
	"github.com/google/uuid"
)

type Route struct {
	ID         string
	Path       string
	TemplateID string
	Meta       struct {
		Title       string
		Description string
	}
	Blocks []blocks.Handle
}

type Router struct {
	Store map[string]Route
}

func New() Router {
	return Router{
		Store: map[string]Route{},
	}
}

func (router Router) NewRoute() (Route, error) {
	id := uuid.New().String()

	router.Store[id] = Route{
		ID: id,
	}

	return router.Store[id], nil
}

func (router Router) UpdateRoute(route Route) error {
	if _, ok := router.Store[route.ID]; !ok {
		return errors.New("cannot find route")
	}

	for _, r := range router.Store {
		if r.Path == route.Path && r.ID != route.ID {
			return errors.New("cannot update route")
		}
	}

	router.Store[route.ID] = route

	return nil
}

func (router Router) GetRoute(pathOrID string) (Route, error) {
	if route, ok := router.Store[pathOrID]; ok {
		return route, nil
	}

	for _, route := range router.Store {
		if route.Path == pathOrID {
			return route, nil
		}
	}

	return Route{}, errors.New("route not found")
}

func (router Router) DeleteRoute(route Route) error {
	delete(router.Store, route.ID)

	return nil
}
