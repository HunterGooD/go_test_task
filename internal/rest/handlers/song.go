package handlers

import (
	"context"
	"log/slog"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

//go:generate mockery --name SongUsecase
type SongUsecase interface {
	GetListSong(ctx context.Context, filters map[string]string) ([]entity.Song, error)
}

type SongHandler struct {
	usecase SongUsecase
	logger  *slog.Logger
}

func NewSongHandler(r *gin.Engine, usecase SongUsecase, logger *slog.Logger) {
	handler := &SongHandler{usecase, logger}

	r.GET("/songs", handler.GetSongs)
}

func (s *SongHandler) GetSongs(c *gin.Context) {

}
