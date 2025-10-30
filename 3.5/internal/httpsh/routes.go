package httpsh

import (
	"net/http"
	"strings"
	"threeFive/domain"
)

type Route struct {
	Method        string
	Path          string
	Handler       http.HandlerFunc
	MiddlewareLog func(http.Handler) http.Handler
}

func RegisterRoutes(mux *http.ServeMux, handlers domain.Handlers) {
	routes := []Route{
		{Method: http.MethodPost, Path: "/events", Handler: handlers.Events},
		{Method: http.MethodGet, Path: "/events/", Handler: handlers.Get},
	}

	for _, route := range routes {
		handler := applyMiddleware(http.Handler(route.Handler), route.MiddlewareLog)
		mux.Handle(route.Path, methodHandler(route.Method, handler))
	}
	mux.Handle("/events/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/events/")
		parts := strings.Split(strings.Trim(path, "/"), "/")

		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}

		id := parts[0]

		switch {
		case len(parts) == 1 && r.Method == http.MethodGet:
			r.URL.Query().Add("id", id)
			handlers.Get(w, r)
			return

		case len(parts) == 2 && parts[1] == "book" && r.Method == http.MethodPost:
			r.URL.Query().Add("id", id)
			handlers.Book(w, r)
			return

		case len(parts) == 2 && parts[1] == "confirm" && r.Method == http.MethodPost:
			r.URL.Query().Add("id", id)
			handlers.Confirm(w, r)
			return

		default:
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		}
	}))
}

func methodHandler(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func applyMiddleware(h http.Handler, middleware func(http.Handler) http.Handler) http.Handler {
	if middleware != nil {
		return middleware(h)
	}
	return h
}
