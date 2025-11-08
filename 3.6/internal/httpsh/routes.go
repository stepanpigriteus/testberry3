package httpsh

import (
	"net/http"
	"threeSixth/domain"
)

type contextKey string

const eventIDKey contextKey = "eventID"

func RegisterRoutes(mux *http.ServeMux, handlers domain.Handlers) {
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {

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
