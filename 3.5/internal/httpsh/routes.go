package httpsh

import (
	"context"
	"net/http"
	"strings"

	"threeFive/domain"
)

type contextKey string

const eventIDKey contextKey = "eventID"

func RegisterRoutes(mux *http.ServeMux, handlers domain.Handlers) {
	mux.Handle("/events", methodHandler(http.MethodPost, handlers.Create))
	mux.Handle("/user", methodHandler(http.MethodPost, handlers.CreateUser))
	mux.HandleFunc("/events/", eventsRouter(handlers))
}

func eventsRouter(handlers domain.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/events/")
		parts := strings.Split(strings.Trim(path, "/"), "/")

		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, `{"error":"empty path not found"}`, http.StatusNotFound)
			return
		}

		eventID := parts[0]
		ctx := context.WithValue(r.Context(), eventIDKey, eventID)
		r = r.WithContext(ctx)

		// GET /events/{id} — получение информации о мероприятии
		if len(parts) == 1 {
			if r.Method != http.MethodGet {
				http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
				return
			}
			handlers.Gets(w, r)
			return
		}

		// Проверка что есть action
		if len(parts) != 2 {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}

		// POST /events/{id}/book или /events/{id}/confirm
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		action := parts[1]
		switch action {
		case "book":
			handlers.Book(w, r)
		case "confirm":
			handlers.Confirm(w, r)
		default:
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		}
	}
}

func methodHandler(method string, handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

func GetEventID(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(eventIDKey).(string)
	return id, ok
}
