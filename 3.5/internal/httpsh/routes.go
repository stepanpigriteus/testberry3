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

	mux.Handle("/events", methodHandler(http.MethodPost, handlers.Events))
	// Все роуты вида /events/{id}/*
	mux.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/events/")
		parts := strings.SplitN(strings.Trim(path, "/"), "/", 2)

		// GET /events/ - список всех мероприятий
		if len(parts) == 0 || parts[0] == "" {
			if r.Method == http.MethodGet {
				handlers.Gets(w, r)
				return
			}
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		eventID := parts[0]
		ctx := context.WithValue(r.Context(), eventIDKey, eventID)

		// GET /events/{id} - информация о конкретном мероприятии
		if len(parts) == 1 {
			if r.Method == http.MethodGet {
				handlers.Gets(w, r.WithContext(ctx))
				return
			}
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		// Роуты с действиями: /events/{id}/{action}
		action := parts[1]
		switch action {
		case "book":
			if r.Method == http.MethodPost {
				handlers.Book(w, r.WithContext(ctx))
				return
			}
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)

		case "confirm":
			if r.Method == http.MethodPost {
				handlers.Confirm(w, r.WithContext(ctx))
				return
			}
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)

		default:
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		}
	})
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
