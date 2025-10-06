package domain

import "time"

type Notify struct {
	Id       int       `json:"id"`
	Timing   time.Time `json:"timing"`
	Descript string    `json:"descript"`
}
