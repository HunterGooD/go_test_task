package entity

import "time"

type Song struct {
	ID          int64      `json:"id,omitempty"           db:"id"`
	Name        string     `json:"name,omitempty"         db:"m_name"`
	Link        string     `json:"link,omitempty"         db:"m_link"`
	Text        string     `json:"text,omitempty"         db:"m_text"`
	ReleaseDate time.Time  `json:"release_date"           db:"m_release_date"`
	CreatedAt   time.Time  `json:"created_at,omitempty"   db:"created_at"`
	UpdateAt    time.Time  `json:"update_at,omitempty"    db:"update_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"   db:"deleted_at"`
	GroupID     int64      `json:"group_id,omitempty"     db:"group_id"`
	Group       *Group     `json:"group,omitempty"        db:"group"`
}

type SongRequest struct {
	Group       string     `json:"group" validate:"required"`
	Song        string     `json:"song"  validate:"required"`
	Link        string     `json:"link,omitempty"`
	Text        string     `json:"text,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
}

type SongListResponse struct {
	Total   int    `json:"total"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Songs   []Song `json:"songs"`
}

type SongTextResponse struct {
	TotalPages int      `json:"total_pages"`
	Page       int      `json:"page"`
	Text       []string `json:"text"`
}

// SongListQueryParams struct for bind query param ?id=1&name=asda
type SongListQueryParams struct {
	ID          int64      `json:"id,omitempty"           form:"id,omitempty"`
	Name        string     `json:"name,omitempty"         form:"name,omitempty"`
	Link        string     `json:"link,omitempty"         form:"link,omitempty"`
	Text        string     `json:"text,omitempty"         form:"text,omitempty"`
	Page        int        `json:"p,omitempty"            form:"p,omitempty"`
	Limit       int        `json:"limit,omitempty"        form:"limit,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty" form:"release_date,omitempty"`
	GroupName   string     `json:"group_name,omitempty"   form:"group_name,omitempty"`
}

// SongFilters struct for bind body json if text biggest
type SongFilters struct {
	ID          int64      `json:"id,omitempty"           db:"id"`
	Name        string     `json:"name,omitempty"         db:"m_name"`
	Link        string     `json:"link,omitempty"         db:"m_link"`
	Text        string     `json:"text,omitempty"         db:"m_text"`
	ReleaseDate *time.Time `json:"release_date,omitempty" db:"m_release_date"`
	GroupName   string     `json:"group_name,omitempty"   db:"g_name"`
}
