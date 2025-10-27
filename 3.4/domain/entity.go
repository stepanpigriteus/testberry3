package domain

import (
	"image"
	"time"
)

type ImageData struct {
	Bytes       []byte
	ContentType string
	Filename    string
}

type ImageStat struct {
	ID            string    `db:"id" json:"id"`
	OriginalPath  string    `db:"original_path" json:"original_path"`
	ThumbnailPath *string   `db:"thumbnail_path,omitempty" json:"thumbnail_path,omitempty"`
	WatermarkPath *string   `db:"watermark_path,omitempty" json:"watermark_path,omitempty"`
	ResizedPath   *string   `db:"resized_path,omitempty" json:"resized_path,omitempty"`
	Status        string    `db:"status" json:"status"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type Task struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Bucket   string `json:"bucket"`
}

type Outputs struct {
	Thumbnail image.Image
	Watermark image.Image
	Resized   image.Image
}
