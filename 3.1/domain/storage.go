package domain

type Storage interface {
	CreateNotify(notify Notify) error
	GetNotify(id string) (error, Notify)
	DeleteNotify(id string) error
}
