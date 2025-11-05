package httpsh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info().Msg("Events handler called")
	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "не удалось прочитать JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	id, err := h.serv.Create(ctx, event)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidSeats):
			writeError(w, http.StatusBadRequest, fmt.Sprintf("неверное количество мест: %v", err))
		case errors.Is(err, domain.ErrAlreadyExists):
			writeError(w, http.StatusBadRequest, fmt.Sprintf("событие уже существует: %v", err))
		default:
			writeError(w, http.StatusBadRequest, fmt.Sprintf("внутренняя ошибка сервера: %v", err))
		}

	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok", "eventId": id})
}

func (h *Handlers) Gets(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("EventGets handler called")
	ctx := r.Context()
	// /events/{id}
	eventId := strings.TrimPrefix(r.URL.Path, "/events/")
	event, err := h.serv.Gets(ctx, eventId)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("не удалось получить сщобытие: %v", err))
	}
	writeJSON(w, http.StatusOK, event)

}

func (h *Handlers) Book(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Book handler called")
	ctx := r.Context()
	eventID, ok := ctx.Value(eventIDKey).(string)
	if !ok || eventID == "" {
		http.Error(w, "eventId not found in context", http.StatusBadRequest)
		return
	}
	userID := r.URL.Query().Get("user")
	if userID == "" {
		http.Error(w, "userId not found in query", http.StatusBadRequest)
		return
	}
	bookID, err := h.serv.Book(ctx, eventID, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create booking")

		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.Error(w, "event not found", http.StatusNotFound)
		case errors.Is(err, domain.ErrAlreadyExists):
			http.Error(w, "booking already exists", http.StatusConflict)
		case errors.Is(err, domain.ErrInvalidBooking):
			http.Error(w, "invalid booking data", http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"status":  "ok",
		"book_id": bookID,
	})
}

func (h *Handlers) Confirm(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Confirm handler called")
	ctx := r.Context()
	eventID, ok := ctx.Value(eventIDKey).(string)
	if !ok || eventID == "" {
		http.Error(w, "eventId not found in context", http.StatusBadRequest)
		return
	}
	bookId := r.URL.Query().Get("book")
	if bookId == "" {
		writeError(w, http.StatusBadRequest, "не указан id брони ")
	}
	err := h.serv.Confirm(ctx, eventID, bookId)
	if err != nil {
		h.logger.Err(err).Msg("не удалось подвтердить бронь")
		writeError(w, http.StatusBadRequest, fmt.Sprintf("не удалось подвтердить бронь: %v", err))
	}

}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("CreateUser handler called")
	ctx := r.Context()
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	id, err := h.serv.CreateUser(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		h.logger.Error().Err(err).Msg("failed to create user")
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{
		"id": id,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
