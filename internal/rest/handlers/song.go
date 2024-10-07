package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/utils"
	"github.com/gin-gonic/gin"
)

const DEFAULT_LIMIT_SONG int = 10

//go:generate mockery --name SongUsecase
type SongUsecase interface {
	CreateNewSong(ctx context.Context, songInput *entity.SongRequest) (*entity.Song, error)
	GetListSong(ctx context.Context, page, pageSize int, filters *entity.SongFilters) ([]entity.Song, error)
	TotalSongs(ctx context.Context) (int, error)
	DeleteSoftByID(ctx context.Context, id int64) error
	DeleteForceByID(ctx context.Context, id int64) error
}

//go:generate mockery --name MusicInfoUsecase
type MusicInfoUsecase interface {
	GetInfo(ctx context.Context, songInput *entity.SongRequest) error
}

type SongHandler struct {
	songUsecase      SongUsecase
	musicInfoUsecase MusicInfoUsecase
}

func NewSongHandler(r *gin.Engine, usecase SongUsecase, musicInfo MusicInfoUsecase) {
	handler := &SongHandler{usecase, musicInfo}

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
	var page, limit int
	result := &entity.SongListResponse{}
	querySong := &entity.SongListQueryParams{}
	filtersSong := &entity.SongFilters{}

	ctx := c.Request.Context()

	if err := c.BindQuery(querySong); err != nil {
		if !errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, entity.ErrorResponse{
				Code:    400,
				Error:   entity.ErrBadParamInput.Error(),
				Message: "Error parse query param",
			})
			slog.Error("cannot bind query", slog.String("error", err.Error()))
			return
		}
		err = nil // to nil because skipping error if io.EOF
	}

	slog.Info("get query param", slog.Any("querySong", querySong), slog.String("url", c.Request.URL.RawQuery))

	if err := c.BindJSON(filtersSong); err != nil {
		if !errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, entity.ErrorResponse{
				Code:    400,
				Error:   entity.ErrBadParamInput.Error(),
				Message: "Error parse filter body",
			})
			slog.Error("cannot bind json", slog.String("error", err.Error()))
			return
		}
		err = nil // to nil because skipping error if io.EOF
	}

	limit, page = querySong.Limit, querySong.Page

	if querySong.Limit == 0 {
		limit = DEFAULT_LIMIT_SONG
	}

	if querySong.Page == 0 {
		page = 1
	}

	filtersSong = utils.MergeSongParams(querySong, filtersSong)

	listSongs, err := s.songUsecase.GetListSong(ctx, page, limit, filtersSong)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			slog.Error("error on function `s.songUsecase.GetListSong` with params",
				slog.String("error", err.Error()),
				slog.Int("page", page),
				slog.Int("limit", limit),
				slog.Any("filtesSong", filtersSong),
			)
			c.JSON(http.StatusNotFound, entity.ErrorResponse{
				Code:    404,
				Error:   entity.ErrNotFound.Error(),
				Message: "Error usecase get list",
			})
			return
		}
		slog.Error("error on function `s.songUsecase.GetListSong` with params",
			slog.String("error", err.Error()),
			slog.Int("page", page),
			slog.Int("limit", limit),
			slog.Any("filtesSong", filtersSong),
		)
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:    400,
			Error:   entity.ErrBadParamInput.Error(),
			Message: "Error usecase get list",
		})
		return
	}

	total, err := s.songUsecase.TotalSongs(ctx)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			slog.Error("error on function `s.songUsecase.TotalSongs`",
				slog.String("error", err.Error()),
			)
			c.JSON(http.StatusNotFound, entity.ErrorResponse{
				Code:    404,
				Error:   entity.ErrNotFound.Error(),
				Message: "Error usecase total songs",
			})
			return
		}
		slog.Error("error on function `s.songUsecase.TotalSongs`",
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:    400,
			Error:   entity.ErrBadParamInput.Error(),
			Message: "Error usecase total songs",
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
// @Param song_insert body entity.SongRequest true "Song insert"
// @Success 200 {object} entity.Song           "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/create [post]
func (s *SongHandler) CreateNewSong(c *gin.Context) {
	var songInput *entity.SongRequest
	ctx := c.Request.Context()

	err := s.musicInfoUsecase.GetInfo(ctx, songInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	err = c.BindJSON(songInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:  400,
			Error: entity.ErrBadParamInput.Error(),
		})
		return
	}

	songRes, err := s.songUsecase.CreateNewSong(ctx, songInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:  400,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, songRes)
}

// @Summary Delete song
// @Schemes
// @Description Deleting song
// @Tags Song
// @Accept json
// @Produce json
// @Param song_id path  int  true  "Song id"
// @Param soft    query bool false "Is soft delete"
// @Success 200 "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/{song_id} [delete]
func (s *SongHandler) DeleteSong(c *gin.Context) {
	// songID, err := strconv.Atoi(c.Param("song_id")) return int but id is int64
	isSoft := c.Query("soft")
	ctx := c.Request.Context()
	songID, err := strconv.ParseInt(c.Param("song_id"), 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:  400,
			Error: entity.ErrBadParamInput.Error(),
		})
		return
	}
	if len(isSoft) != 0 {
		// soft deleting
		err = s.songUsecase.DeleteSoftByID(ctx, songID)
	} else {
		// force delete
		err = s.songUsecase.DeleteForceByID(ctx, songID)
	}
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, entity.ErrorResponse{
				Code:  404,
				Error: entity.ErrNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:  400,
			Error: entity.ErrBadParamInput.Error(),
		})
		return
	}
	c.String(http.StatusOK, "ok")
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
