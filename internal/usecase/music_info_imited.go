package usecase

import (
	"context"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
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
	log interfaces.Logger
}

func NewMusicInfoUsecaseImited(logger interfaces.Logger) *MusicInfoUsecaseImited {
	return &MusicInfoUsecaseImited{logger}
}

func (m *MusicInfoUsecaseImited) GetInfo(ctx context.Context, songInput *entity.SongRequest) error {
	m.log.Info("Music info mock searching", map[string]any{
		"song_name":  songInput.Song,
		"group_name": songInput.Group,
	})
	for _, v := range FIXED_VALUES {
		if v.Song == songInput.Song && v.Group == songInput.Group {
			songInput.Link = v.Link
			songInput.Text = v.Text
			songInput.ReleaseDate = v.ReleaseDate
			return nil
		}
	}
	m.log.Warn("error music info not found", map[string]any{
		"code":  404,
		"song":  songInput.Song,
		"group": songInput.Group,
	})
	return entity.ErrNotFound
}
