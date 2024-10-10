package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
	"github.com/HunterGooD/go_test_task/pkg/utils"
	"github.com/gin-gonic/gin"
)

const DEFAULT_LIMIT_SONG int = 10

//go:generate mockery --name SongUsecase
type SongUsecase interface {
	CreateNewSong(ctx context.Context, songInput *entity.SongRequest) (*entity.Song, error)
	GetTextSong(ctx context.Context, songID int64) (*entity.SongTextResponse, error)
	GetListSong(ctx context.Context, page, pageSize int, isDeleted bool, filters *entity.SongFilters) (*entity.SongListResponse, error)
	FullUpdateSong(ctx context.Context, song *entity.Song) error
	UpdateSong(ctx context.Context, song *entity.Song) error
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
	log              interfaces.Logger
}

func NewSongHandler(r *gin.Engine, usecase SongUsecase, musicInfo MusicInfoUsecase, logger interfaces.Logger) {
	logger.Info("create song handler")
	handler := &SongHandler{usecase, musicInfo, logger}

	// Getting songs handlers
	r.GET("/song/list", handler.GetSongs)
	logger.Info("SONG register handler path GET `/song/list`")

	r.GET("/song/:song_id/text", handler.GetText)
	logger.Info("SONG register handler path GET `/song/:song_id/text`")

	// Create handlers
	r.POST("/song/create", handler.CreateNewSong)
	logger.Info("SONG register handler path POST `/song/create`")

	// Delete handlers
	r.DELETE("/song/:song_id", handler.DeleteSong)
	logger.Info("SONG register handler path DELETE `/song/:song_id`")

	// Updates handlers
	r.PUT("/song/:song_id", handler.PutSong)
	logger.Info("SONG register handler path PUT `/song/:song_id`")

	r.PATCH("/song/:song_id", handler.PatchSong)
	logger.Info("SONG register handler path PATCH `/song/:song_id`")
}

// @Summary get list songs with filters
// @Schemes
// @Description Getting songs with pagination and filtered
// @Tags Song
// @Accept json
// @Produce json
// @Param   d            query    bool                        false  "With deleted"
// @Param   query_params query    entity.SongListQueryParams  false  "Query params"
// @Param   filters      body     entity.SongFilters          false  "Filters json body"
// @Success 200 {object} entity.SongListResponse "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Failure 500 {object} entity.ErrorResponse "Internal server error"
// @Router /song/list [get]
func (s *SongHandler) GetSongs(c *gin.Context) {
	var page, limit int
	result := &entity.SongListResponse{}
	querySong := &entity.SongListQueryParams{}
	filtersSong := &entity.SongFilters{}

	isDelete := true
	with_deleted, err := strconv.ParseBool(c.Query("d"))
	if err != nil {
		s.log.Warn("Parse bool with deleted set default value `false`")
		isDelete = false
	} else {
		s.log.Info("set value with deleted", map[string]any{
			"with_deleted": with_deleted,
		})
		isDelete = with_deleted
	}

	ctx := c.Request.Context()

	if err := c.ShouldBindQuery(querySong); err != nil {
		if err != io.EOF {
			c.JSON(http.StatusBadRequest, entity.ErrorResponse{
				Code:    400,
				Error:   err.Error(),
				Message: entity.ErrBadParamInput.Error(),
			})
			s.log.Error("cannot bind query filter", map[string]any{
				"err": err,
			})
			return
		}
		s.log.Warn("no query filter")
	}

	s.log.Info("get query param", map[string]any{
		"querySong": querySong,
		"url_query": c.Request.URL.RawQuery,
	})

	if err := c.ShouldBindJSON(filtersSong); err != nil {
		if err != io.EOF {
			c.JSON(http.StatusBadRequest, entity.ErrorResponse{
				Code:    400,
				Error:   err.Error(),
				Message: entity.ErrBadParamInput.Error(),
			})
			s.log.Error("cannot bind body json filter", map[string]any{
				"err": err,
			})
			return
		}
		s.log.Warn("no query filter")
	}

	limit, page = querySong.Limit, querySong.Page

	if querySong.Limit == 0 {
		limit = DEFAULT_LIMIT_SONG
	}

	if querySong.Page == 0 {
		page = 1
	}

	filtersSong = utils.MergeSongParams(querySong, filtersSong)
	s.log.Debug("final filter", map[string]any{
		"filter": filtersSong,
	})

	s.log.Info("get list song")
	result, err = s.songUsecase.GetListSong(ctx, page, limit, isDelete, filtersSong)
	if err != nil {
		s.log.Error("error on function `s.songUsecase.GetListSong` with params", map[string]any{
			"error":      err,
			"page":       page,
			"limit":      limit,
			"filtesSong": filtersSong,
		})
		if err == entity.ErrNotFound {
			c.JSON(http.StatusNotFound, entity.ErrorResponse{
				Code:    404,
				Error:   entity.ErrNotFound.Error(),
				Message: fmt.Sprintf("Not found with filters=%#v", filtersSong),
			})
		}
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:    400,
			Error:   err.Error(),
			Message: entity.ErrBadParamInput.Error(),
		})
		return
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
	ctx := c.Request.Context()
	songID, err := strconv.ParseInt(c.Param("song_id"), 10, 0)
	if err != nil {
		s.log.Error("error get text song", map[string]any{
			"err": err,
		})
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:    400,
			Error:   err.Error(),
			Message: entity.ErrBadParamInput.Error(),
		})
		return
	}

	text, err := s.songUsecase.GetTextSong(ctx, songID)
	if err != nil {
		s.log.Error("error get text song", map[string]any{
			"err": err,
		})
		if err == entity.ErrNotFound {
			c.JSON(http.StatusNotFound, entity.ErrorResponse{
				Code:    http.StatusNotFound,
				Error:   err.Error(),
				Message: entity.ErrNotFound.Error(),
			})
		}
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, text)
}

// @Summary Create song
// @Schemes
// @Description Creating song with song name and group name
// @Tags Song
// @Accept json
// @Produce json
// @Param song_insert body entity.SongRequest true "Song insert"
// @Success 200 {object} entity.Song          "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/create [post]
func (s *SongHandler) CreateNewSong(c *gin.Context) {
	songInput := &entity.SongRequest{}
	ctx := c.Request.Context()

	err := c.BindJSON(songInput)
	if err != nil {
		s.log.Error("error bind json body", map[string]any{
			"err": err,
		})
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:    400,
			Message: entity.ErrBadParamInput.Error(),
			Error:   err.Error(),
		})
		return
	}

	err = s.musicInfoUsecase.GetInfo(ctx, songInput)
	if err != nil {
		s.log.Error("error from music info service", map[string]any{
			"err": err,
		})
		if err == entity.ErrNotFound {
			c.JSON(http.StatusNotFound, entity.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Error on `s.musicInfoUsecase.GetInfo` songInput=%v", songInput),
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error on `s.musicInfoUsecase.GetInfo` songInput=%v", songInput),
			Error:   err.Error(),
		})
		return
	}

	songRes, err := s.songUsecase.CreateNewSong(ctx, songInput)
	if err != nil {
		s.log.Error("error from music info service", map[string]any{
			"err": err,
		})
		if err == entity.ErrNotFound {
			c.JSON(http.StatusNotFound, entity.ErrorResponse{
				Code:  http.StatusNotFound,
				Error: err.Error(),
			})
			return
		}
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
// @Description Deleting song if soft delete true then deleted_at in BD change on NOW() or if soft is false then delete row from bd
// @Tags Song
// @Accept json
// @Produce json
// @Param song_id path  int  true  "Song id"
// @Param soft    query bool false "Is soft delete"
// @Success 200  "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/{song_id} [delete]
func (s *SongHandler) DeleteSong(c *gin.Context) {
	// songID, err := strconv.Atoi(c.Param("song_id")) return int but id is int64
	ctx := c.Request.Context()
	songID, err := strconv.ParseInt(c.Param("song_id"), 10, 0)
	if err != nil {
		s.log.Error("parse song_id error", map[string]any{
			"err": err,
		})
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:  400,
			Error: entity.ErrBadParamInput.Error(),
		})
		return
	}

	isSoft := false
	parseSoft, err := strconv.ParseBool(c.Query("soft"))
	if err != nil {
		isSoft = true
		s.log.Warn("error parse soft query param setting default true")
	} else {
		isSoft = parseSoft
	}
	if isSoft {
		// soft deleting
		err = s.songUsecase.DeleteSoftByID(ctx, songID)
	} else {
		// force delete
		err = s.songUsecase.DeleteForceByID(ctx, songID)
	}
	if err != nil {
		s.log.Error("error on delete ", map[string]any{
			"err": err,
		})
		if err == entity.ErrNotFound {
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
	c.Status(http.StatusOK)
}

// @Summary Put song
// @Schemes
// @Description Put song change all fields
// @Tags Song
// @Accept json
// @Produce json
// @Param song_id path int         true "Song id"
// @Param song    body entity.Song true "Song with changing fields"
// @Success 200 {object} entity.Song           "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/{song_id} [put]
func (s *SongHandler) PutSong(c *gin.Context) {
	song := &entity.Song{}
	ctx := c.Request.Context()
	songID, err := strconv.ParseInt(c.Param("song_id"), 10, 0)
	if err != nil {
		s.log.Error("error on parse id", map[string]any{
			"err": err,
		})
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:  400,
			Error: entity.ErrBadParamInput.Error(),
		})
		return
	}
	if err := c.BindJSON(song); err != nil {
		s.log.Error("error on bind song", map[string]any{
			"err": err,
		})
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:  400,
			Error: entity.ErrBadParamInput.Error(),
		})
		return
	}

	// set id on struct for update and where query
	if song.ID == 0 {
		song.ID = songID
	}

	if err := s.songUsecase.FullUpdateSong(ctx, song); err != nil { // check on 404
		s.log.Error("error from full update song", map[string]any{
			"err":  err,
			"song": song,
		})
		if err == entity.ErrNotFound {
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

	c.JSON(http.StatusOK, song)
}

// @Summary Patch song
// @Schemes
// @Description Patch song change any fields
// @Tags Song
// @Accept json
// @Produce json
// @Param song_id path int         true "Song id"
// @Param song    body entity.Song true "Song with changing fields"
// @Success 200 {object} entity.Song           "ok"
// @Failure 400 {object} entity.ErrorResponse "Params not valid"
// @Failure 404 {object} entity.ErrorResponse "Can not find ID"
// @Router /song/{song_id} [patch]
func (s *SongHandler) PatchSong(c *gin.Context) {
	song := &entity.Song{}
	ctx := c.Request.Context()
	songID, err := strconv.ParseInt(c.Param("song_id"), 10, 0)
	if err != nil {
		s.log.Error("error parse id", map[string]any{
			"err": err,
		})
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:  400,
			Error: entity.ErrBadParamInput.Error(),
		})
		return
	}

	if err := c.BindJSON(song); err != nil {
		s.log.Error("error on bind song", map[string]any{
			"err": err,
		})
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:  400,
			Error: entity.ErrBadParamInput.Error(),
		})
		return
	}

	if song.ID == 0 {
		song.ID = songID
	}

	if err := s.songUsecase.UpdateSong(ctx, song); err != nil {
		s.log.Error("error from update song", map[string]any{
			"err":  err,
			"song": song,
		})
		if err == entity.ErrNotFound {
			c.JSON(http.StatusNotFound, entity.ErrorResponse{
				Code:  404,
				Error: entity.ErrNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Code:    400,
			Error:   err.Error(),
			Message: entity.ErrBadParamInput.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, song)
}
