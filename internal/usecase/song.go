package usecase

import (
	"context"

	"github.com/HunterGooD/go_test_task/internal/entity"
)

//go:generate mockery --name SongRepository
type SongRepository interface {
	GetSongs(ctx context.Context, filters map[string]string) ([]entity.Song, error)
}

type GroupRepository interface {
}

type SongUsecase struct {
	songRepo  SongRepository
	groupRepo GroupRepository
}

func NewSongUsecase(sr SongRepository, gr GroupRepository) *SongUsecase {
	return &SongUsecase{sr, gr}
}

func (su *SongUsecase) GetSongs(ctx context.Context, filters map[string]string) ([]entity.Song, error) {
	res, err := su.songRepo.GetSongs(ctx, filters)
	if err != nil {
		return nil, err
	}

	return res, err
}
