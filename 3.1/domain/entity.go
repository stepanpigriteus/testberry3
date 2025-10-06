package domain

import "time"

type Notify struct {
	Id        int       `json:"id"`
	Timing    time.Time `json:"timing"`
	Descript  string    `json:"descript"`
	Status    string    `json:"status"`
	Retry     int       `json:"retry"`
	CreatedAt time.Time `json:"created_at"`
}
