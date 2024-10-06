package utils_test

import (
	"testing"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestRef(t *testing.T) {
	tables := []struct {
		name     string
		a        entity.SongFilters
		expected string
	}{
		{
			"Test with id column",
			entity.SongFilters{
				ID: 1,
			},
			"AND id = 1",
		},
		{
			"Test with id and group_name column",
			entity.SongFilters{
				ID:        1,
				GroupName: "GroupNameasdasd",
			},
			"AND id = 1 AND g_name LIKE 'GroupNameasdasd%'",
		},
		{
			"Test with all column",
			entity.SongFilters{
				ID:          1,
				Name:        "Nameasda",
				Link:        "Linkasd",
				Text:        "Textasda",
				ReleaseDate: time.Date(2024, 10, 6, 13, 13, 10, 0, time.Now().UTC().Location()),
				GroupName:   "GroupNameasdasd",
			},
			"AND id = 1 AND m_name LIKE 'Nameasda%' AND m_link LIKE 'Linkasd%' AND m_text LIKE 'Textasda%' AND m_release_date = '2024-10-06' AND g_name LIKE 'GroupNameasdasd%'",
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			res := utils.GetFilterString(&tt.a)
			assert.Equal(t, tt.expected, res)
		})
	}
}
