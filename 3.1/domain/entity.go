package domain

import "time"

type Notify struct {
	id       string
	timing   time.Time
	descript string
}
