package entity

import "time"

type Song struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Link        string    `json:"link"`
	Text        string    `json:"text" validate:"required"`
	ReleaseDate time.Time `json:"release_date"`
	Group       Group     `json:"group"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
