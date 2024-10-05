package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type GroupUsecase interface {
}

type GroupHandler struct {
	usecase GroupUsecase
	logger  *slog.Logger
}

func NewGroupHandler(r *gin.Engine, usecase GroupUsecase, logger *slog.Logger) {
	handler := &GroupHandler{usecase, logger}

	r.GET("/groups", handler.GetGroups)
}

func (g *GroupHandler) GetGroups(c *gin.Context) {

}
