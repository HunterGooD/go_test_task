package usecase

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/pkg/api" // TODO: Interface and entity depends without strong dependency
)

//go:generate mockery --name MusicInfoApi
type MusicInfoApi interface {
	GetInfo(ctx context.Context, params *api.GetInfoParams, reqEditors ...api.RequestEditorFn) (*http.Response, error)
}

type MusicInfoUsecase struct {
	musicInfoApi MusicInfoApi
}

func NewMusicInfoUsecase(musicInfoApi MusicInfoApi) *MusicInfoUsecase {
	return &MusicInfoUsecase{musicInfoApi}
}

func (m *MusicInfoUsecase) GetInfo(ctx context.Context, songInput *entity.SongRequest) error {
	response, err := m.musicInfoApi.GetInfo(ctx, &api.GetInfoParams{
		Group: songInput.Group,
		Song:  songInput.Song,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()
	// Проверка статуса ответа
	if response.StatusCode != http.StatusOK {
		return entity.ErrNotFound
	}

	// Чтение тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return entity.ErrBadParamInput
	}

	err = json.Unmarshal(body, &songInput)
	if err != nil {
		return entity.ErrBadParamInput
	}

	return nil
}
