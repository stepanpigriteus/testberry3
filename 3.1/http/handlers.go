package http

import (
	"net/http"
	"treeOne/domain"

	"github.com/rs/zerolog"
)

type HandleNotify struct {
	logger  zerolog.Logger
	service domain.Service
}

func NewHandleNotify(logger zerolog.Logger, service domain.Service) *HandleNotify {
	return &HandleNotify{
		logger:  logger,
		service: service,
	}
}

func (h *HandleNotify) CreateNotify(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Endpoint works")
}

func (h *HandleNotify) GetNotify(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Endpoint works")
}

func (h *HandleNotify) DeleteNotify(w http.ResponseWriter, r *http.Request) {

}
