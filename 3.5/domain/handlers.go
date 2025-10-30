package domain

import "net/http"

type Handlers interface {
	Gets(w http.ResponseWriter, r *http.Request)
	Book(w http.ResponseWriter, r *http.Request)
	Confirm(w http.ResponseWriter, r *http.Request)
	Events(w http.ResponseWriter, r *http.Request)
}
