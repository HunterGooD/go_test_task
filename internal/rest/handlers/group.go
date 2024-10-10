package handlers

import (
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type GroupUsecase interface {
}

type GroupHandler struct {
	usecase GroupUsecase
	log     interfaces.Logger
}

func NewGroupHandler(r *gin.Engine, usecase GroupUsecase, logger interfaces.Logger) {
	logger.Info("Create new group handler")
	handler := &GroupHandler{usecase, logger}

	r.GET("/group/list", handler.GetGroups)
	logger.Info("GROUP register handler path GET `/group/list`")

}

func (g *GroupHandler) GetGroups(c *gin.Context) {

}
