package utils_test

import (
	"testing"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestRef(t *testing.T) {
	d := time.Date(2024, 10, 6, 13, 13, 10, 0, time.Now().UTC().Location())
	tables := []struct {
		name        string
		a           *entity.SongFilters
		expected    string
		expectedLen int
	}{
		{
			"Test with id column",
			&entity.SongFilters{
				ID: 1,
			},
			"AND id = $3",
			1,
		},
		{
			"Test with id and group_name column",
			&entity.SongFilters{
				ID:        1,
				GroupName: "GroupNameasdasd",
			},
			"AND id = $3 AND g_name LIKE $4",
			2,
		},
		{
			"Test with all column",
			&entity.SongFilters{
				ID:          1,
				Name:        "Nameasda",
				Link:        "Linkasd",
				Text:        "Textasda",
				ReleaseDate: &d,
				GroupName:   "GroupNameasdasd",
			},
			"AND id = $3 AND m_name LIKE $4 AND m_link LIKE $5 AND m_text LIKE $6 AND m_release_date = $7 AND g_name LIKE $8",
			6,
		},
		{
			"Test with nil",
			nil,
			"",
			0,
		},
		{
			"Test with empty struct",
			&entity.SongFilters{},
			"",
			0,
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			res, args := utils.GetFilterString(3, tt.a)
			assert.Equal(t, tt.expected, res)
			assert.Len(t, args, tt.expectedLen)
		})
	}
}

func TestMergeSongParams(t *testing.T) {
	d := time.Date(2024, 10, 6, 13, 13, 10, 0, time.Now().UTC().Location())
	tables := []struct {
		name     string
		a        *entity.SongFilters
		b        *entity.SongListQueryParams
		expected *entity.SongFilters
	}{
		{
			"Test with one filter param",
			&entity.SongFilters{
				ID: 1,
			},
			&entity.SongListQueryParams{
				ID:   2,
				Name: "qweq",
			},
			&entity.SongFilters{
				ID:   1,
				Name: "qweq",
			},
		},
		{
			"Test with all empty filter param",
			&entity.SongFilters{},
			&entity.SongListQueryParams{
				ID:          2,
				Name:        "qweq",
				Link:        "123",
				Text:        "213",
				ReleaseDate: &d,
				GroupName:   "qweq",
			},
			&entity.SongFilters{
				ID:          2,
				Name:        "qweq",
				Link:        "123",
				Text:        "213",
				ReleaseDate: &d,
				GroupName:   "qweq",
			},
		},
		{
			"Test with all fill filter param",
			&entity.SongFilters{
				ID:          1,
				Name:        "qweq",
				Link:        "123",
				Text:        "213",
				ReleaseDate: &d,
				GroupName:   "qweq",
			},
			&entity.SongListQueryParams{
				ID:          2,
				Name:        "2",
				Link:        "2",
				Text:        "2",
				ReleaseDate: &d,
				GroupName:   "2",
			},
			&entity.SongFilters{
				ID:          1,
				Name:        "qweq",
				Link:        "123",
				Text:        "213",
				ReleaseDate: &d,
				GroupName:   "qweq",
			},
		},
		{
			"Test with nil",
			nil,
			nil,
			&entity.SongFilters{},
		},
	}
	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			res := utils.MergeSongParams(tt.b, tt.a)
			assert.Equal(t, tt.expected, res)
		})
	}
}
