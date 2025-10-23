package domain

import "time"

type ImageData struct {
	Bytes       []byte
	ContentType string
	Filename    string
}

type ImageStat struct {
	ID            string    `db:"id" json:"id"`
	OriginalPath  string    `db:"original_path" json:"original_path"`
	ProcessedPath *string   `db:"processed_path,omitempty" json:"processed_path,omitempty"`
	Status        string    `db:"status" json:"status"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
