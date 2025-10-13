package domain

import "net/http"

type Handlers interface {
	CreateShorten(w http.ResponseWriter, r *http.Request)
	GetShorten(w http.ResponseWriter, r *http.Request)
	GetAnalytics(w http.ResponseWriter, r *http.Request)
}
