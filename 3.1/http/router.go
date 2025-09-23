package http

import (
	"fmt"
	"net/http"
	"treeOne/domain"
)

type Route struct {
	Method        string
	Path          string
	Handler       http.HandlerFunc
	MiddlewareLog func(http.Handler) http.Handler
}

func RegisterRoutes(mux *http.ServeMux, handlers domain.EventHandler) {
	routes := []Route{}
	for _, route := range routes {
		finalHandler := applyMiddleware(route.Handler, route.MiddlewareLog)
		mux.Handle(route.Path, methodHandler(route.Method, finalHandler))
	}
}

func methodHandler(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(">>> methodHandler: got %s request on %s\n", r.Method, r.URL.Path)
		if r.Method != method {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func applyMiddleware(h http.Handler, middleware func(http.Handler) http.Handler) http.Handler {
	h = middleware(h)
	return h
}
