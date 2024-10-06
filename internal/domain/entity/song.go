package entity

import "time"

type Song struct {
	ID          int64      `json:"id,omitempty"           db:"id"`
	Name        string     `json:"name"                   db:"m_name"`
	Link        string     `json:"link"                   db:"m_link"`
	Text        string     `json:"text"                   db:"m_text"`
	ReleaseDate time.Time  `json:"release_date"           db:"m_release_date"`
	CreatedAt   time.Time  `json:"created_at,omitempty"   db:"created_at"`
	UpdateAt    time.Time  `json:"update_at,omitempty"    db:"update_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"   db:"deleted_at"`
	GroupID     int64      `json:"-"                      db:"group_id"`
	Group       *Group     `json:"group,omitempty"        db:"group"`
}

type SongInsert struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
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
	ID          int64     `json:"id"           from:"id"`
	Name        string    `json:"name"         from:"name"`
	Link        string    `json:"link"         from:"link"`
	Text        string    `json:"text"         from:"text"`
	Page        int       `json:"p"            from:"p"`
	Limit       int       `json:"limit"        from:"limit"`
	ReleaseDate time.Time `json:"release_date" from:"release_date"`
	GroupName   string    `json:"group_name"   from:"group_name"`
}

// SongFilters struct for bind body json if text biggest
type SongFilters struct {
	ID          int64     `json:"id,omitempty"           db:"id"`
	Name        string    `json:"name,omitempty"         db:"m_name"`
	Link        string    `json:"link,omitempty"         db:"m_link"`
	Text        string    `json:"text,omitempty"         db:"m_text"`
	ReleaseDate time.Time `json:"release_date,omitempty" db:"m_release_date"`
	GroupName   string    `json:"group_name,omitempty"   db:"g_name"`
}
