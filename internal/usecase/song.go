package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
)

type SongUsecase struct {
	songRepo                      interfaces.SongRepository
	transactionManagerSongsGroups interfaces.TransactionManagerSongsGroups
	log                           interfaces.Logger
}

// NewSongUsecase operations with songs get, change, delete
func NewSongUsecase(sr interfaces.SongRepository, tmSG interfaces.TransactionManagerSongsGroups, logger interfaces.Logger) *SongUsecase {
	return &SongUsecase{sr, tmSG, logger}
}

func (su *SongUsecase) CreateNewSong(ctx context.Context, songInput *entity.SongRequest) (*entity.Song, error) {
	su.log.Info("starting create song", map[string]any{
		"song":  songInput.Song,
		"group": songInput.Group,
	})
	// Start transaction for add song
	err := su.transactionManagerSongsGroups.Begin()
	if err != nil {
		su.log.Error("begin transaction", map[string]any{
			"err":   err,
			"song":  songInput.Song,
			"group": songInput.Group,
		})
		return nil, err
	}

	su.log.Info("begin trasaction")
	su.log.Debug("starting create song", map[string]any{
		"song":  songInput.Song,
		"group": songInput.Group,
	})

	defer func() {
		if err != nil {
			su.log.Error("rollback transaction", map[string]any{
				"err":        err,
				"song_input": songInput,
			})
			if err := su.transactionManagerSongsGroups.Rollback(); err != nil {
				su.log.Error("rollback transaction not finish", map[string]any{
					"err":        err,
					"song_input": songInput,
				})
			}
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
		su.log.Error("error on get group by name", map[string]any{
			"err":        err,
			"song_input": songInput,
		})
		group, err = txGroupRepo.CreateGroup(ctx, songInput.Group)
		if err != nil {
			su.log.Error("error on create group", map[string]any{
				"err":        err,
				"song_input": songInput,
			})
			return nil, err
		}

	}
	if group == nil {
		su.log.Error("error getting group", map[string]any{
			"song_input": songInput,
		})
		return nil, entity.ErrNotFound
	}

	su.log.Info("getting group", map[string]any{
		"group_id": group.ID,
	})
	su.log.Debug("getting group", map[string]any{
		"group": group,
	})

	// check if not exists song and create song
	song, err = txSongRepo.GetByName(ctx, songInput.Song, group.ID)
	if err != nil {
		su.log.Error("error on get song by name", map[string]any{
			"err":        err,
			"song_input": songInput,
		})
		song, err = txSongRepo.CreateSong(ctx, group.ID, songInput)
		if err != nil {
			su.log.Error("error on create song by name", map[string]any{
				"err":        err,
				"song_input": songInput,
			})
			return nil, err
		}
	}

	if song == nil {
		su.log.Error("error getting song", map[string]any{
			"song_input": songInput,
		})
		return nil, entity.ErrNotFound
	}

	su.log.Info("getting song", map[string]any{
		"song_id": song.ID,
	})
	su.log.Debug("getting song", map[string]any{
		"song": song,
	})

	if err = su.transactionManagerSongsGroups.Commit(); err != nil {
		su.log.Error("error transaction commit", map[string]any{
			"err":        err,
			"song":       song,
			"group":      group,
			"song_input": songInput,
		})
		return nil, err
	}

	su.log.Info("succes transaction with commit", map[string]any{
		"song_id":    song.ID,
		"song_name":  song.Name,
		"group_id":   group.ID,
		"group_name": group.GName,
	})

	return song, nil
}

// GetListSong select all songs
func (su *SongUsecase) GetListSong(ctx context.Context, page, pageSize int, isDeleted bool, filters *entity.SongFilters) (*entity.SongListResponse, error) {
	offset := (page - 1) * pageSize

	songs, err := su.songRepo.GetListSong(ctx, offset, pageSize, isDeleted, filters)
	if err != nil {
		su.log.Error("error get list song", map[string]any{
			"err": err,
		})
		return nil, err
	}

	total, err := su.songRepo.Total(ctx, isDeleted, filters)
	if err != nil {
		su.log.Error("error get total songs", map[string]any{
			"err": err,
		})
		return nil, err
	}

	su.log.Info("gettings song", map[string]any{
		"total":        total,
		"length_songs": len(songs),
	})

	// if page > max page return reverse songs in limit
	if total != 0 && len(songs) == 0 {
		songs, err = su.songRepo.GetReverseListSongs(ctx, 0, pageSize, isDeleted, filters)
		if err != nil {
			su.log.Error("error get songs reverse", map[string]any{
				"err": err,
			})
			return nil, err
		}
		su.log.Info("gettings reverse songs", map[string]any{
			"total":        total,
			"length_songs": len(songs),
		})
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
		su.log.Error("error get text song", map[string]any{
			"err": err,
		})
		return nil, err
	}

	su.log.Info("text getting and in process pagination...", map[string]any{
		"song_id":     songID,
		"text_length": len(text),
	})

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

	su.log.Info("text pagination complited", map[string]any{
		"song_id":     songID,
		"text_length": len(text),
		"pages":       res.TotalPages,
	})
	su.log.Debug("text pagination complited", map[string]any{
		"song_id":         songID,
		"text_length":     len(text),
		"text_pagination": res,
	})

	return res, nil
}

func (su *SongUsecase) FullUpdateSong(ctx context.Context, song *entity.Song) error {
	var params map[string]any

	if song == nil {
		params = map[string]any{
			"m_name":         "",
			"m_link":         "",
			"m_text":         "",
			"m_release_date": time.Now(),
		}
	} else {
		params = map[string]any{
			"m_name":         song.Name,
			"m_link":         song.Link,
			"m_text":         song.Text,
			"m_release_date": song.ReleaseDate,
		}
	}

	su.log.Info("update song", map[string]any{
		"song_id": song.ID,
	})

	err := su.songRepo.UpdateFromMapByID(ctx, song.ID, song, params)
	if err != nil {
		su.log.Error("error update song", map[string]any{
			"err": err,
		})
		return err
	}
	su.log.Debug("updated song", map[string]any{
		"song_id": song.ID,
		"song":    song,
	})
	return nil
}

func (su *SongUsecase) UpdateSong(ctx context.Context, song *entity.Song) error {
	params := make(map[string]any)
	if song == nil {
		return nil
	}
	if len(song.Name) != 0 {
		params["m_name"] = song.Name
	}
	if len(song.Link) != 0 {
		params["m_link"] = song.Link
	}
	if len(song.Text) != 0 {
		params["m_text"] = song.Text
	}
	if song.DeletedAt != nil {
		params["deleted_at"] = song.DeletedAt
	}
	if !song.ReleaseDate.IsZero() {
		params["m_release_date"] = song.ReleaseDate
	}

	if len(params) == 0 {
		return nil
	}

	su.log.Info("update song", map[string]any{
		"song_id":       song.ID,
		"length_params": len(params),
	})

	su.log.Debug("update song", map[string]any{
		"song_id":       song.ID,
		"length_params": len(params),
		"params":        params,
	})

	err := su.songRepo.UpdateFromMapByID(ctx, song.ID, song, params)
	if err != nil {
		su.log.Error("error update song", map[string]any{
			"err": err,
		})
		return err
	}

	su.log.Debug("updated song", map[string]any{
		"song_id": song.ID,
		"song":    song,
	})
	return nil

}

func (su *SongUsecase) GetSongTextByID(ctx context.Context, id int64) {

}

func (su *SongUsecase) DeleteSoftByID(ctx context.Context, id int64) error {
	su.log.Info("start soft delete", map[string]any{
		"song_id": id,
	})
	return su.songRepo.DeleteSoftByID(ctx, id)
}

func (su *SongUsecase) DeleteForceByID(ctx context.Context, id int64) error {
	su.log.Info("start force delete", map[string]any{
		"song_id": id,
	})
	return su.songRepo.DeleteForceByID(ctx, id)
}
