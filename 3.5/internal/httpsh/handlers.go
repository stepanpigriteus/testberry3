package httpsh

import (
	"context"
	"net/http"

	"threeFive/domain"

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

func (h *Handlers) Events(w http.ResponseWriter, r *http.Request) {
}

func (h *Handlers) Gets(w http.ResponseWriter, r *http.Request) {
}

func (h *Handlers) Book(w http.ResponseWriter, r *http.Request) {
}

func (h *Handlers) Confirm(w http.ResponseWriter, r *http.Request) {
}
