package handlers

import (
	"github.com/gin-gonic/gin"
)

type GroupUsecase interface {
}

type GroupHandler struct {
	usecase GroupUsecase
}

func NewGroupHandler(r *gin.Engine, usecase GroupUsecase) {
	handler := &GroupHandler{usecase}

	r.GET("/group", handler.GetGroups)
}

func (g *GroupHandler) GetGroups(c *gin.Context) {

}
