package domain

type EventHandler interface {
	CreateEvent()
	GetEventStatus()
	DeleteEvent()
}
