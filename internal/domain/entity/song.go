package entity

import "time"

type Song struct {
	ID          int64      `json:"id"           db:"id"`
	Name        string     `json:"name"         db:"m_name" validate:"required"`
	Link        string     `json:"link"         db:"m_link"`
	Text        string     `json:"text"         db:"m_text" validate:"required"`
	ReleaseDate time.Time  `json:"release_date" db:"m_release_date"`
	CreatedAt   time.Time  `json:"created_at"   db:"created_at"`
	UpdateAt    time.Time  `json:"update_at"    db:"update_at"`
	DeletedAt   *time.Time `json:"deleted_at"   db:"deleted_at"`
	GroupID     int64      `json:"group_id"     db:"group_id"`
	Group       Group      `json:"group"        db:"group"`
}
