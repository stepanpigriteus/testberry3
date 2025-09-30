package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"treeOne/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/redis"
)

type HandleNotify struct {
	logger  zerolog.Logger
	service domain.Service
	redis   redis.Client
}

func NewHandleNotify(logger zerolog.Logger, service domain.Service, redis redis.Client) *HandleNotify {
	return &HandleNotify{
		logger:  logger,
		service: service,
		redis:   redis,
	}
}

func (h *HandleNotify) CreateNotify(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Endpoint works")
	ctx := r.Context()
	var notify domain.Notify
	if err := json.NewDecoder(r.Body).Decode(&notify); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.CreateNotify(ctx, notify)
	if err != nil {
		http.Error(w, "ошибка проброшена в хендлер: "+err.Error(), http.StatusBadRequest)
		return
	}

}

func (h *HandleNotify) GetNotify(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Endpoint works")
	ctx := context.Background()
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	id := parts[2]
	notify, err := h.service.GetNotify(ctx, id)
	if err != nil {
		http.Error(w, "ошибка проброшена в хендлер: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(notify)
}

func (h *HandleNotify) DeleteNotify(w http.ResponseWriter, r *http.Request) {

}
