package httpsh

import (
	"context"
	"net/http"
	"threeSixth/domain"

	"github.com/rs/zerolog"
)

type Handlers struct {
	serv   domain.Service
	logger zerolog.Logger
}

func NewHandlers(ctx context.Context, serv domain.Service, logger zerolog.Logger) *Handlers {
	return &Handlers{
		serv:   serv,
		logger: logger,
	}
}

func (h *Handlers) CreateItem(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetItems(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetItem(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) UpdateItem(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) DeleteItem(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetAnalytics(w http.ResponseWriter, r *http.Request) {}
