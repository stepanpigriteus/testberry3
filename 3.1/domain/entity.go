package domain

import "time"

type Notify struct {
	Id       int       `json:"Id"`
	Timing   time.Time `json:"Timing"`
	Descript string    `json:"Descript"`
}
