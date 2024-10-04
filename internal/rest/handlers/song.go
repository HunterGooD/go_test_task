package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type SongUsecase interface {
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
