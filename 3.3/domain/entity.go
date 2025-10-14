package domain

import "time"

type Comment struct {
	ID        int64      `db:"id" json:"id"`
	ParentID  *int64     `db:"parent_id" json:"parent_id,omitempty"`
	Text      string     `db:"text" json:"text"`
	Author    string     `db:"author" json:"author"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Children  []*Comment `json:"children,omitempty"`
}
