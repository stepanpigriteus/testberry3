package domain

type EventHandler interface {
	CreateNotify()
	GetNotify()
	DeleteNotify()
}
