package entity

import "time"

type Group struct {
	ID        int64      `json:"id,omitempty"         db:"id"`
	GName     string     `json:"g_name"               db:"g_name" validate:"required"`
	CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdateAt  time.Time  `json:"updated_at,omitempty" db:"update_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"` // DeletedAt nil value if not deleted
}
