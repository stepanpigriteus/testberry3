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
		writeError(w, http.StatusBadRequest, "не удалось прочитать JSON: "+err.Error())
		return
	}
	defer r.Body.Close()

	id, err := h.serv.Create(ctx, event)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidSeats):
			writeError(w, http.StatusBadRequest, "неверное количество мест")
		case errors.Is(err, domain.ErrAlreadyExists), errors.Is(err, domain.ErrDuplicateKey):
			writeError(w, http.StatusConflict, "событие уже существует")
		default:
			writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		}
		return
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
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "событие не найдено")
			return
		}
		writeError(w, http.StatusInternalServerError, "не удалось получить событие")
		return
	}
	writeJSON(w, http.StatusOK, event)
}

func (h *Handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("GetAll handler called")
	ctx := r.Context()
	events, err := h.serv.GetAll(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось получить события")
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (h *Handlers) Book(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Book handler called")
	ctx := r.Context()
	eventID, ok := ctx.Value(eventIDKey).(string)
	if !ok || eventID == "" {
		writeError(w, http.StatusBadRequest, "eventId not found in context")
		return
	}
	userID := r.URL.Query().Get("user")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "userId not found in query")
		return
	}
	bookID, err := h.serv.Book(ctx, eventID, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create booking")

		switch {
		case errors.Is(err, domain.ErrNotFound):
			writeError(w, http.StatusNotFound, "событие не найдено")
		case errors.Is(err, domain.ErrUserNotFound):
			writeError(w, http.StatusNotFound, "пользователь не найден")
		case errors.Is(err, domain.ErrAlreadyExists), errors.Is(err, domain.ErrDuplicateKey):
			writeError(w, http.StatusConflict, "бронь уже существует")
		case errors.Is(err, domain.ErrInvalidSeats):
			writeError(w, http.StatusBadRequest, "нет доступных мест")
		case errors.Is(err, domain.ErrInvalidBooking):
			writeError(w, http.StatusBadRequest, "некорректные данные брони")
		default:
			writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
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
		writeError(w, http.StatusBadRequest, "eventId not found in context")
		return
	}
	bookId := r.URL.Query().Get("book")
	if bookId == "" {
		writeError(w, http.StatusBadRequest, "не указан id брони")
		return
	}
	err := h.serv.Confirm(ctx, eventID, bookId)
	if err != nil {
		h.logger.Error().Err(err).Msg("не удалось подтвердить бронь")

		switch {
		case errors.Is(err, domain.ErrNotFound):
			writeError(w, http.StatusNotFound, "бронь не найдена")
		case errors.Is(err, domain.ErrInvalidStatus):
			writeError(w, http.StatusBadRequest, "невозможно подтвердить бронь с текущим статусом")
		default:
			writeError(w, http.StatusInternalServerError, "не удалось подтвердить бронь")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("CreateUser handler called")
	ctx := r.Context()
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	fmt.Println(user)
	id, err := h.serv.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateKey) {
			writeError(w, http.StatusConflict, "пользователь с таким email уже существует")
			return
		}
		if strings.Contains(err.Error(), "required") {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.logger.Error().Err(err).Msg("failed to create user")
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{
		"id": id,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
