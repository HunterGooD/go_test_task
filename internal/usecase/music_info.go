package usecase

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
	"github.com/HunterGooD/go_test_task/pkg/api" // TODO: Interface and entity depends without strong dependency
)

//go:generate mockery --name MusicInfoApi
type MusicInfoApi interface {
	GetInfo(ctx context.Context, params *api.GetInfoParams, reqEditors ...api.RequestEditorFn) (*http.Response, error)
}

type MusicInfoUsecase struct {
	musicInfoApi MusicInfoApi

	log interfaces.Logger
}

func NewMusicInfoUsecase(musicInfoApi MusicInfoApi, logger interfaces.Logger) *MusicInfoUsecase {
	logger.Info("creat music info usecase")
	return &MusicInfoUsecase{musicInfoApi, logger}
}

func (m *MusicInfoUsecase) GetInfo(ctx context.Context, songInput *entity.SongRequest) error {
	m.log.Info("start get info client", map[string]any{
		"song":  songInput.Song,
		"group": songInput.Group,
	})

	response, err := m.musicInfoApi.GetInfo(ctx, &api.GetInfoParams{
		Group: songInput.Group,
		Song:  songInput.Song,
	})
	if err != nil {
		m.log.Error("error client request", map[string]any{
			"err":   err,
			"song":  songInput.Song,
			"group": songInput.Group,
		})
		return err
	}

	defer response.Body.Close()
	// Проверка статуса ответа
	if response.StatusCode != http.StatusOK {
		m.log.Warn("error music info not found", map[string]any{
			"code":  response.StatusCode,
			"song":  songInput.Song,
			"group": songInput.Group,
		})
		return entity.ErrNotFound
	}

	// Чтение тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return entity.ErrBadParamInput
	}

	err = json.Unmarshal(body, &songInput)
	if err != nil {
		m.log.Error("error music info json body parse", map[string]any{
			"song_name":   songInput.Song,
			"group_name":  songInput.Group,
			"json_string": string(body),
			"song":        songInput,
		})
		return entity.ErrBadParamInput
	}

	m.log.Info("getting music info success")

	return nil
}
