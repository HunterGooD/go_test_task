package usecase

import (
	"context"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
)

var releaseDate time.Time = time.Now()

var FIXED_VALUES []entity.SongRequest = []entity.SongRequest{
	{
		Group:       "TestGroup1",
		Song:        "TestSong1",
		Link:        "TestLink1",
		Text:        "Test test test test\nTest test test test\nTest test test test\nTest test test test\n\nTest test test test\nTest test test test\nTest test test test\nTest test test test\n\nTest test test test\nTest test test test\nTest test test test\nTest test test test",
		ReleaseDate: &releaseDate,
	},
	{
		Group:       "TestGroup2",
		Song:        "TestSong2",
		Link:        "TestLink2",
		Text:        "Test test test test\nTest test test test\nTest test test test\nTest test test test\n\nTest test test test\nTest test test test\nTest test test test\nTest test test test\n\nTest test test test\nTest test test test\nTest test test test\nTest test test test",
		ReleaseDate: &releaseDate,
	},
	{
		Group:       "TestGroup3",
		Song:        "TestSong3",
		Link:        "TestLink3",
		Text:        "Test test test test\nTest test test test\nTest test test test\nTest test test test\n\n\nTest test test test\nTest test test test\nTest test test test\nTest test test test\n\n\nTest test test test\nTest test test test\nTest test test test\nTest test test test",
		ReleaseDate: &releaseDate,
	},
	{
		Group:       "TestGroup4",
		Song:        "TestSong4",
		Link:        "TestLink4",
		Text:        "Test test test test\nTest test test test\nTest test test test\nTest test test test\n\n\n\nTest test test test\nTest test test test\nTest test test test\nTest test test test\n\n\n\nTest test test test\nTest test test test\nTest test test test\nTest test test test",
		ReleaseDate: &releaseDate,
	},
}

type MusicInfoUsecaseImited struct {
}

func NewMusicInfoUsecaseImited() *MusicInfoUsecaseImited {
	return &MusicInfoUsecaseImited{}
}

func (m *MusicInfoUsecaseImited) GetInfo(ctx context.Context, songInput *entity.SongRequest) error {
	for _, v := range FIXED_VALUES {
		if v.Song == songInput.Song && v.Group == songInput.Group {
			songInput.Link = v.Link
			songInput.Text = v.Text
			songInput.ReleaseDate = v.ReleaseDate
			return nil
		}
	}
	return entity.ErrNotFound
}
