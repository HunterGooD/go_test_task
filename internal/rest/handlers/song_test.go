package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/rest/handlers"
	"github.com/HunterGooD/go_test_task/internal/rest/handlers/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test for GET /song/list
func TestGetSongs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	type SUTotal struct {
		returnObject int
		returnError  error
	}
	type SUGetListSong struct {
		returnObject []entity.Song
		returnError  error
	}
	tables := []struct {
		nameTest            string
		reqURL              string
		mockSUGetListSong   SUGetListSong
		mockSUTotal         SUTotal
		queryParams         *entity.SongListQueryParams
		filterBody          *entity.SongFilters
		expectedCode        int
		expectedErrorStruct *entity.ErrorResponse
	}{
		{
			nameTest: "test succes with correct params",
			reqURL:   "/song/list",
			mockSUGetListSong: SUGetListSong{
				returnObject: []entity.Song{
					{
						ID:   1,
						Name: "song 1",
					},
					{
						ID:   2,
						Name: "song 2",
					},
				},
				returnError: nil,
			},
			mockSUTotal: SUTotal{
				returnObject: 2,
				returnError:  nil,
			},
			queryParams: &entity.SongListQueryParams{
				Page:  1,
				Limit: 10,
			},
			filterBody:          &entity.SongFilters{},
			expectedCode:        200,
			expectedErrorStruct: nil,
		},
		{
			nameTest: "test error without params",
			reqURL:   "/song/list?p=11111&limit=9999999",
			mockSUGetListSong: SUGetListSong{
				returnObject: []entity.Song{},
				returnError:  nil,
			},
			mockSUTotal: SUTotal{
				returnObject: 0,
				returnError:  entity.ErrBadParamInput,
			},
			queryParams:  &entity.SongListQueryParams{Page: 11111, Limit: 9999999},
			filterBody:   &entity.SongFilters{},
			expectedCode: 400,
			expectedErrorStruct: &entity.ErrorResponse{
				Code:    400,
				Message: "Error usecase total songs",
				Error:   "params is not valid",
			},
		},
	}

	for _, tt := range tables {
		t.Run(tt.nameTest, func(t *testing.T) {
			router := gin.Default()
			mockSongusecase := new(mocks.SongUsecase)
			mockMusicInfoUsecase := new(mocks.MusicInfoUsecase)
			querySong := tt.queryParams
			filterSong := tt.filterBody
			returnMockSongUsecase := tt.mockSUGetListSong.returnObject
			totalSong := tt.mockSUTotal.returnObject
			expectedRes := &entity.SongListResponse{}

			if querySong.Limit == 0 {
				querySong.Limit = handlers.DEFAULT_LIMIT_SONG
			}

			if querySong.Page == 0 {
				querySong.Page = 1
			}
			expectedRes = &entity.SongListResponse{
				Total:   totalSong,
				Page:    querySong.Page,
				PerPage: querySong.Limit,
				Songs:   returnMockSongUsecase,
			}

			filterBodyJSON, _ := json.Marshal(filterSong)
			var expectedInJSON []byte
			if tt.expectedErrorStruct != nil {
				expectedInJSON, _ = json.Marshal(tt.expectedErrorStruct)

			} else {
				expectedInJSON, _ = json.Marshal(expectedRes)
			}

			mockSongusecase.On("GetListSong", mock.Anything, querySong.Page, querySong.Limit, filterSong).
				Return(returnMockSongUsecase, tt.mockSUGetListSong.returnError).
				Once()

			mockSongusecase.On("TotalSongs", mock.Anything).
				Return(totalSong, tt.mockSUTotal.returnError).
				Once()

			handlers.NewSongHandler(router, mockSongusecase, mockMusicInfoUsecase)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", tt.reqURL, strings.NewReader(string(filterBodyJSON)))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.NotEmpty(t, string(expectedInJSON))

			if tt.expectedErrorStruct == nil {
				result := &entity.SongListResponse{}
				err := json.Unmarshal(w.Body.Bytes(), result)
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, expectedRes, result)
			} else {
				resultErr := &entity.ErrorResponse{}
				err := json.Unmarshal(w.Body.Bytes(), resultErr)
				assert.NoError(t, err)
				assert.NotNil(t, resultErr)
				assert.Equal(t, tt.expectedErrorStruct, resultErr)
			}

		})
	}

}
