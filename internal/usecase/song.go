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

// GetListSong select all songs
func (su *SongUsecase) GetListSong(ctx context.Context, page, pageSize int, filters *entity.SongFilters) ([]entity.Song, error) {
	offset := (page - 1) * pageSize
	return su.songRepo.GetListSong(ctx, offset, pageSize, filters)
}

func (su *SongUsecase) TotalSongs(ctx context.Context) (int, error) {
	return su.songRepo.Total(ctx)
}

func (su *SongUsecase) GetSongTextByID(ctx context.Context, id int64) {

}
