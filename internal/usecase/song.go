package usecase

import (
	"context"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
)

type SongUsecase struct {
	songRepo                      interfaces.SongRepository
	transactionManagerSongsGroups interfaces.TransactionManagerSongsGroups
}

// NewSongUsecase operations with songs get, change, delete
func NewSongUsecase(sr interfaces.SongRepository, tmSG interfaces.TransactionManagerSongsGroups) *SongUsecase {
	return &SongUsecase{sr, tmSG}
}

func (su *SongUsecase) GetListSong(ctx context.Context, filters map[string]string) ([]entity.Song, error) {
	res, err := su.songRepo.GetListSong(ctx, 10, 0, filters)
	if err != nil {
		return nil, err
	}

	return res, err
}
