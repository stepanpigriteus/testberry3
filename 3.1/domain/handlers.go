package domain

import "net/http"

type EventHandler interface {
	CreateNotify(w http.ResponseWriter, r *http.Request)
	GetNotify(w http.ResponseWriter, r *http.Request)
	DeleteNotify(w http.ResponseWriter, r *http.Request)
}
