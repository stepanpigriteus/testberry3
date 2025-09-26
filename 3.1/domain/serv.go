package domain

type Service interface {
	CreateNotify(notify Notify) error
	GetNotify(id string) (error, Notify)
	DeleteNotify(id string) error
}
