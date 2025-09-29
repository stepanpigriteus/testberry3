package domain

import "time"

type Notify struct {
	Id       string
	Timing   time.Time
	Descript string
}
