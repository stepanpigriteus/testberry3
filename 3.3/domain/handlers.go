package domain

import "net/http"

type Handlers interface {
	CreateComments(w http.ResponseWriter, r *http.Request)
	GetComments(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
