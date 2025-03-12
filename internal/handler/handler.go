package handler

import (
	"errors"
	"sync"

	"github.com/Fluffy-Bean/cms/internal/blocks"
	"github.com/google/uuid"
)

type Page struct {
	ID         string
	Path       string
	TemplateID string
	Meta       struct {
		Title       string
		Description string
	}
	Blocks []blocks.Handle
}

type Handler struct {
	mu sync.Mutex

	TemplatesPath string
	DataPath      string
	Pages         map[string]Page
	Blocks        blocks.Blocks
}

func (h *Handler) NewPage() (Page, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	id := uuid.New().String()

	h.Pages[id] = Page{
		ID: id,
	}

	return h.Pages[id], nil
}

func (h *Handler) UpdatePage(page Page) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.Pages[page.ID]; !ok {
		return errors.New("cannot find page")
	}

	for _, r := range h.Pages {
		if r.Path == page.Path && r.ID != page.ID {
			return errors.New("cannot update page")
		}
	}

	h.Pages[page.ID] = page

	return nil
}

func (h *Handler) GetPage(pathOrID string) (Page, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if page, ok := h.Pages[pathOrID]; ok {
		return page, nil
	}

	for _, page := range h.Pages {
		if page.Path == pathOrID {
			return page, nil
		}
	}

	return Page{}, errors.New("page not found")
}

func (h *Handler) DeletePage(page Page) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.Pages, page.ID)

	return nil
}
