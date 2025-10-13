package httpsh

import (
	"context"
	"net/http"
	"treethree/domain"

	"github.com/rs/zerolog"
)

type CommentsHandlers struct {
	serv   domain.Service
	logger zerolog.Logger
}

func NewHandlers(ctx context.Context, serv domain.Service, logger zerolog.Logger) *CommentsHandlers {
	return &CommentsHandlers{
		serv:   serv,
		logger: logger,
	}
}

func (h *CommentsHandlers) CreateComments(w http.ResponseWriter, r *http.Request) {

}
func (h *CommentsHandlers) GetComments(w http.ResponseWriter, r *http.Request) {

}

func (h *CommentsHandlers) Delete(w http.ResponseWriter, r *http.Request) {

}
