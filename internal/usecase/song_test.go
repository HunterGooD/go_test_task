package usecase_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces/mocks"
	"github.com/HunterGooD/go_test_task/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

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
		songRepo.On("GetListSong", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("bool"), mock.AnythingOfType("*entity.SongFilters")).
			Return(mockReturnRepo, nil).Once()
		songRepo.On("Total", mock.Anything, mock.AnythingOfType("bool"), mock.AnythingOfType("*entity.SongFilters")).
			Return(2, nil).Once()
		songUsecase := usecase.NewSongUsecase(songRepo, transactionManager, logger)
		songList, err := songUsecase.GetListSong(context.TODO(), 1, 10, false, nil)
		assert.NotEmpty(t, songList)
		assert.NoError(t, err)
		assert.Len(t, songList.Songs, len(mockReturnRepo))

		songRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		songRepo.On("GetListSong", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("bool"), mock.AnythingOfType("*entity.SongFilters")).
			Return(nil, entity.ErrNotFound).Once()

		songUsecase := usecase.NewSongUsecase(songRepo, transactionManager, logger)
		songList, err := songUsecase.GetListSong(context.TODO(), 1, 10, false, nil)
		assert.Nil(t, songList)
		assert.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrNotFound)

		songRepo.AssertExpectations(t)
	})
}

func TestCreateSong(t *testing.T) {
	t.Run("Test success", func(t *testing.T) {
		mockSongRepo := new(mocks.SongRepository)
		mockTransactionManager := new(mocks.TransactionManagerSongsGroups)
		mockGroupRepo := new(mocks.GroupRepository)
		songInput := &entity.SongRequest{
			Song:  "Test Song",
			Group: "Test Group",
		}
		mockTransactionManager.On("Begin").Return(nil)
		mockTransactionManager.On("SongRepository").Return(mockSongRepo)
		mockTransactionManager.On("GroupRepository").Return(mockGroupRepo)
		mockTransactionManager.On("Commit").Return(nil)

		group := &entity.Group{ID: 1, GName: "Test Group"}
		mockGroupRepo.On("GetByName", mock.Anything, "Test Group").Return(group, nil)

		songExcepted := &entity.Song{ID: 1, Name: "Test Song"}
		mockSongRepo.On("GetByName", mock.Anything, "Test Song", group.ID).Return(nil, entity.ErrNotFound)
		mockSongRepo.On("CreateSong", mock.Anything, group.ID, mock.Anything).Return(songExcepted, nil)

		songUsecase := usecase.NewSongUsecase(mockSongRepo, mockTransactionManager, logger)

		song, err := songUsecase.CreateNewSong(context.TODO(), songInput)

		assert.NoError(t, err)
		assert.NotNil(t, song)
		assert.Equal(t, songExcepted, song)
	})

	t.Run("Test error", func(t *testing.T) {
		mockSongRepo := new(mocks.SongRepository)
		mockTransactionManager := new(mocks.TransactionManagerSongsGroups)
		mockGroupRepo := new(mocks.GroupRepository)
		songInput := &entity.SongRequest{
			Song:  "Test Song",
			Group: "Test Group",
		}
		mockTransactionManager.On("Begin").Return(nil)
		mockTransactionManager.On("SongRepository").Return(mockSongRepo)
		mockTransactionManager.On("GroupRepository").Return(mockGroupRepo)
		mockTransactionManager.On("Rollback").Return(nil)

		mockGroupRepo.On("GetByName", mock.Anything, "Test Group").Return(nil, entity.ErrNotFound)
		mockGroupRepo.On("CreateGroup", mock.Anything, "Test Group").Return(nil, entity.ErrNotFound)

		songUsecase := usecase.NewSongUsecase(mockSongRepo, mockTransactionManager, logger)

		song, err := songUsecase.CreateNewSong(context.TODO(), songInput)

		assert.Error(t, err)
		assert.Nil(t, song)
	})

}
