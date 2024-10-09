package usecase

import (
	"context"
	"log/slog"
	"strings"

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

	slog.Info("begin trasaction")
	slog.Debug("begin trasaction", slog.Any("songInput", songInput))

	defer func() {
		if err != nil {
			su.transactionManagerSongsGroups.Rollback()
			slog.Info("rollback trasaction", slog.Any("error", err))
			slog.Debug("rollback trasaction", slog.Any("songInput", songInput))
		}
	}()

	// create repositories with transaction
	txSongRepo := su.transactionManagerSongsGroups.SongRepository()
	txGroupRepo := su.transactionManagerSongsGroups.GroupRepository()

	//entities
	var group *entity.Group
	var song *entity.Song
	// select or create group
	group, err = txGroupRepo.GetByName(ctx, songInput.Group)
	if err != nil {
		group, err = txGroupRepo.CreateGroup(ctx, songInput.Group)
		if err != nil {
			return nil, err
		}

	}
	if group == nil {
		return nil, entity.ErrNotFound
	}
	slog.Info("getting group in transaction", slog.Int64("group_id", group.ID))
	slog.Debug("getting group in transaction", slog.Any("group", group))

	slog.Info("Try get song ", slog.String("name_song", songInput.Song), slog.Int64("group_id", group.ID))
	// check if not exists song and create song
	song, err = txSongRepo.GetByName(ctx, songInput.Song, group.ID)
	if err != nil {
		slog.Info("Try create song ", slog.String("name_song", songInput.Song), slog.Int64("group_id", group.ID))
		slog.Info("Try create song ", slog.Any("song", songInput), slog.Any("group", group))
		song, err = txSongRepo.CreateSong(ctx, group.ID, songInput)
		if err != nil {
			return nil, err
		}

	}
	if song == nil {
		return nil, entity.ErrNotFound
	}
	slog.Info("creating song in transaction", slog.Int64("song_id", song.ID))
	slog.Debug("creating song in transaction", slog.Any("song", song))

	if err = su.transactionManagerSongsGroups.Commit(); err != nil {
		return nil, err
	}
	slog.Info("succes transaction with commit", slog.Int64("song_id", song.ID), slog.Int64("group_id", group.ID))
	slog.Debug("succes transaction with commit", slog.Any("song", song), slog.Any("group", group))

	return song, err
}

// GetListSong select all songs
func (su *SongUsecase) GetListSong(ctx context.Context, page, pageSize int, isDeleted bool, filters *entity.SongFilters) (*entity.SongListResponse, error) {
	offset := (page - 1) * pageSize
	songs, err := su.songRepo.GetListSong(ctx, offset, pageSize, isDeleted, filters)
	if err != nil {
		return nil, err
	}
	total, err := su.songRepo.Total(ctx, isDeleted, filters)
	if err != nil {
		return nil, err
	}
	slog.Info("getting songs", slog.Int("total", total), slog.Int("length_songs", len(songs)))
	if total != 0 && len(songs) == 0 {
		songs, err = su.songRepo.GetReverseListSongs(ctx, 0, pageSize, isDeleted, filters)
		if err != nil {
			return nil, err
		}
		slog.Info("geting reverse song", slog.Int("total", total), slog.Int("length_songs", len(songs)))
	}
	songResponse := &entity.SongListResponse{
		Total:   total,
		Page:    page,
		PerPage: pageSize,
		Songs:   songs,
	}
	return songResponse, nil
}

func (su *SongUsecase) TotalSongs(ctx context.Context, isDel bool, filter *entity.SongFilters) (int, error) {
	return su.songRepo.Total(ctx, isDel, filter)
}

func (su *SongUsecase) GetTextSong(ctx context.Context, songID int64) (*entity.SongTextResponse, error) {
	text, err := su.songRepo.GetSongTextByID(ctx, songID)
	if err != nil {
		return nil, err
	}
	splited := strings.Split(text, "\n\n") // "\n\n" is verse splitting
	resText := make([]string, 0)
	for _, v := range splited {
		if len(v) != 0 {
			resText = append(resText, v)
		}
	}
	res := &entity.SongTextResponse{
		Page:       1,
		TotalPages: len(resText),
		Text:       resText,
	}

	return res, nil
}

func (su *SongUsecase) GetSongTextByID(ctx context.Context, id int64) {

}

func (su *SongUsecase) DeleteSoftByID(ctx context.Context, id int64) error {
	return su.songRepo.DeleteSoftByID(ctx, id)
}

func (s *SongUsecase) DeleteForceByID(ctx context.Context, id int64) error {
	return nil
}
