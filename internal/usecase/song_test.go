package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces/mocks"
	"github.com/HunterGooD/go_test_task/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetListSong(t *testing.T) {
	songRepo := new(mocks.SongRepository)
	transactionManager := new(mocks.TransactionManagerSongsGroups)
	mockReturnRepo := []entity.Song{
		{
			ID:          1,
			Name:        "name song 1",
			Link:        "hello",
			Text:        "qweasdzxc123",
			ReleaseDate: time.Now(),
			GroupID:     1,
			UpdateAt:    time.Now(),
			CreatedAt:   time.Now(),
			DeletedAt:   nil,
			Group: &entity.Group{
				ID:        1,
				GName:     "friks",
				CreatedAt: time.Now(),
				UpdateAt:  time.Now(),
				DeletedAt: nil,
			}},
		{
			ID:          2,
			Name:        "name song 2",
			Link:        "world",
			Text:        "qweasdzxc123",
			ReleaseDate: time.Now(),
			GroupID:     1,
			UpdateAt:    time.Now(),
			CreatedAt:   time.Now(),
			DeletedAt:   nil,
			Group: &entity.Group{
				ID:        1,
				GName:     "friks",
				CreatedAt: time.Now(),
				UpdateAt:  time.Now(),
				DeletedAt: nil,
			}},
	}
	t.Run("success", func(t *testing.T) {
		songRepo.On("GetListSong", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("*entity.SongFilters")).
			Return(mockReturnRepo, nil).Once()
		songUsecase := usecase.NewSongUsecase(songRepo, transactionManager)
		songList, err := songUsecase.GetListSong(context.TODO(), 1, 10, nil)
		assert.NotEmpty(t, songList)
		assert.NoError(t, err)
		assert.Len(t, songList, len(mockReturnRepo))

		songRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		songRepo.On("GetListSong", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("*entity.SongFilters")).
			Return(nil, entity.ErrNotFound).Once()

		songUsecase := usecase.NewSongUsecase(songRepo, transactionManager)
		songList, err := songUsecase.GetListSong(context.TODO(), 1, 10, nil)
		assert.Empty(t, songList)
		assert.Error(t, err)
		assert.Len(t, songList, 0)
		assert.ErrorIs(t, err, entity.ErrNotFound)

		songRepo.AssertExpectations(t)
	})
}
