package httpsh

import (
	"WarehouseControl/domain"
	"net/http"
)

type Route struct {
	Method        string
	Path          string
	Handler       http.HandlerFunc
	MiddlewareLog func(http.Handler) http.Handler
}

func RegisterRoutes(mux *http.ServeMux, handlers domain.Handlers) {
	routes := []Route{}
	pathHandlers := make(map[string][]Route)
	for _, route := range routes {
		pathHandlers[route.Path] = append(pathHandlers[route.Path], route)
	}
	for path, routeList := range pathHandlers {
		mux.Handle(path, multiMethodHandler(routeList))
	}
}

func multiMethodHandler(routes []Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routes {
			if r.Method == route.Method {
				handler := http.Handler(route.Handler)
				finalHandler := applyMiddleware(handler, route.MiddlewareLog)
				finalHandler.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	})
}

func applyMiddleware(h http.Handler, middleware func(http.Handler) http.Handler) http.Handler {
	if middleware != nil {
		h = middleware(h)
	}
	return h
}
