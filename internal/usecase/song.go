package usecase

import (
	"context"
	"errors"

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

func (su *SongUsecase) CreateNewSong(ctx context.Context, songInput *entity.SongRequest) (*entity.Song, error) {
	// Start transaction for add song
	err := su.transactionManagerSongsGroups.Begin()
	if err != nil {
		return nil, err
	}

	// create repositories with transaction
	txSongRepo := su.transactionManagerSongsGroups.SongRepository()
	txGroupRepo := su.transactionManagerSongsGroups.GroupRepository()

	//entities
	var group *entity.Group
	var song *entity.Song
	// select or create group
	group, err = txGroupRepo.GetByName(ctx, songInput.Group)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			group, err = txGroupRepo.CreateGroup(ctx, songInput.Group)
			if err != nil {
				return nil, err
			}
		}
	}
	if group == nil {
		return nil, entity.ErrNotFound
	}

	// check if not exists song and create song
	song, err = txSongRepo.GetByName(ctx, songInput.Song, group.ID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			song, err = txSongRepo.CreateSong(ctx, group.ID, songInput)
			if err != nil {
				return nil, err
			}
		}
	}
	if song == nil {
		return nil, entity.ErrNotFound
	}

	return song, err
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

func (su *SongUsecase) DeleteSoftByID(ctx context.Context, id int64) error {
	return su.songRepo.DeleteSoftByID(ctx, id)
}

func (s *SongUsecase) DeleteForceByID(ctx context.Context, id int64) error {
	return nil
}
