package domain

import "net/http"

type Handlers interface {
	CreateItem(w http.ResponseWriter, r *http.Request)
	GetItems(w http.ResponseWriter, r *http.Request)
	GetItem(w http.ResponseWriter, r *http.Request)
	UpdateItem(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	GetAnalytics(w http.ResponseWriter, r *http.Request)
}
