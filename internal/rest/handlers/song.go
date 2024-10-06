package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/utils"
	"github.com/gin-gonic/gin"
)

const DEFAULT_LIMIT_SONG int = 10

//go:generate mockery --name SongUsecase
type SongUsecase interface {
	GetListSong(ctx context.Context, page, pageSize int, filters *entity.SongFilters) ([]entity.Song, error)
	TotalSongs(ctx context.Context) (int, error)
}

type SongHandler struct {
	usecase SongUsecase
}

func NewSongHandler(r *gin.Engine, usecase SongUsecase) {
	handler := &SongHandler{usecase}

	r.GET("/song/list", handler.GetSongs)
	r.GET("/song/:song_id/text", handler.GetText)
	r.GET("/song/create", handler.CreateNewSong)

	r.DELETE("/song/:song_id", handler.DeleteSong)

	r.PUT("/song/:song_id", handler.PutSong)
	r.PATCH("/song/:song_id", handler.PatchSong)
}

// @Summary get list songs with filters
// @Schemes
// @Description Getting songs with pagination and filtered
// @Tags Song
// @Accept json
// @Produce json
// @Param   p          query    entity.SongListQueryParams  false  "Page"
// @Param   filters    body     entity.SongFilters          false  "Filters"
// @Success 200 {object} entity.SongListResponse "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/list [get]
func (s *SongHandler) GetSongs(c *gin.Context) {
	var querySong *entity.SongListQueryParams
	var page, limit int
	var result *entity.SongListResponse
	filtersSong := &entity.SongFilters{}

	ctx := c.Request.Context()

	if err := c.BindQuery(querySong); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:    400,
			Message: entity.ErrBadParamInput.Error(),
		})
		return
	}
	if err := c.BindJSON(filtersSong); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:    400,
			Message: entity.ErrBadParamInput.Error(),
		})
		return
	}
	if querySong == nil {
		page, limit = 1, DEFAULT_LIMIT_SONG
	} else {
		if querySong.Limit == 0 {
			limit = querySong.Limit
		}
		if querySong.Page == 0 {
			page = querySong.Page
		}
	}
	filtersSong = utils.MergeSongParams(querySong, filtersSong)

	listSongs, err := s.usecase.GetListSong(ctx, page, limit, filtersSong)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, entity.ErrorResponse{
				Code:    404,
				Message: entity.ErrNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:    400,
			Message: entity.ErrBadParamInput.Error(),
		})
		return
	}

	total, err := s.usecase.TotalSongs(ctx)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, entity.ErrorResponse{
				Code:    404,
				Message: entity.ErrNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:    400,
			Message: entity.ErrBadParamInput.Error(),
		})
		return
	}

	result = &entity.SongListResponse{
		Total:   total,
		Page:    page,
		PerPage: limit,
		Songs:   listSongs,
	}

	c.JSON(http.StatusOK, result)
}

// @Summary get text songs with pagination
// @Schemes
// @Description Get text with pagination
// @Tags Song
// @Accept json
// @Produce json
// @Param song_id path int true "Song id"
// @Success 200 {object} entity.SongTextResponse	"ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/{song_id}/text [get]
func (s *SongHandler) GetText(c *gin.Context) {

}

// @Summary Create song
// @Schemes
// @Description Creating
// @Tags Song
// @Accept json
// @Produce json
// @Param song_insert body entity.SongInsert true "Song insert"
// @Success 200 {string} string               "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/create [post]
func (s *SongHandler) CreateNewSong(c *gin.Context) {

}

// @Summary Delete song
// @Schemes
// @Description Deleting song
// @Tags Song
// @Accept json
// @Produce json
// @Param song_id path int true "Song id"
// @Success 200 {string} string               "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/{song_id} [delete]
func (s *SongHandler) DeleteSong(c *gin.Context) {

}

// @Summary Put song
// @Schemes
// @Description Put song
// @Tags Song
// @Accept json
// @Produce json
// @Param song_id path int true "Song id"
// @Success 200 {string} string               "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/{song_id} [put]
func (s *SongHandler) PutSong(c *gin.Context) {

}

// @Summary Patch song
// @Schemes
// @Description Patch song
// @Tags Song
// @Accept json
// @Produce json
// @Param song_id path int true "Song id"
// @Success 200 {string} string               "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/{song_id} [patch]
func (s *SongHandler) PatchSong(c *gin.Context) {

}
