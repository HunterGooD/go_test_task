package handlers_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/rest/handlers"
	"github.com/HunterGooD/go_test_task/internal/rest/handlers/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func ginTest() {
	gin.SetMode(gin.TestMode)
}

// Test for GET /song/list
func TestGetSongs(t *testing.T) {
	ginTest()

	type SUGetListSong struct {
		returnObject *entity.SongListResponse
		returnError  error
	}
	tables := []struct {
		nameTest            string
		reqURL              string
		mockSUGetListSong   SUGetListSong
		queryParams         *entity.SongListQueryParams
		filterBody          *entity.SongFilters
		expectedCode        int
		expectedErrorStruct *entity.ErrorResponse
	}{
		{
			nameTest: "test succes with correct params",
			reqURL:   "/song/list",
			mockSUGetListSong: SUGetListSong{
				returnObject: &entity.SongListResponse{
					Total:   2,
					Page:    1,
					PerPage: 10,
					Songs: []entity.Song{
						{
							ID:   1,
							Name: "song 1",
						},
						{
							ID:   2,
							Name: "song 2",
						},
					},
				},
				returnError: nil,
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
				returnObject: &entity.SongListResponse{},
				returnError:  entity.ErrBadParamInput,
			},
			queryParams:  &entity.SongListQueryParams{Page: 11111, Limit: 9999999},
			filterBody:   &entity.SongFilters{},
			expectedCode: 400,
			expectedErrorStruct: &entity.ErrorResponse{
				Code:    400,
				Message: "Error usecase get list",
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
			expectedRes := &entity.SongListResponse{}

			if querySong.Limit == 0 {
				querySong.Limit = handlers.DEFAULT_LIMIT_SONG
			}

			if querySong.Page == 0 {
				querySong.Page = 1
			}
			expectedRes = &entity.SongListResponse{
				Total:   len(tt.mockSUGetListSong.returnObject.Songs),
				Page:    querySong.Page,
				PerPage: querySong.Limit,
				Songs:   tt.mockSUGetListSong.returnObject.Songs,
			}

			filterBodyJSON, _ := json.Marshal(filterSong)
			var expectedInJSON []byte
			if tt.expectedErrorStruct != nil {
				expectedInJSON, _ = json.Marshal(tt.expectedErrorStruct)

			} else {
				expectedInJSON, _ = json.Marshal(expectedRes)
			}

			mockSongusecase.On(
				"GetListSong",                               // name function
				mock.Anything,                               // 1 context
				mock.AnythingOfType("int"),                  // 2 page
				mock.AnythingOfType("int"),                  // 3 pageSize
				mock.AnythingOfType("bool"),                 // 4 isDeleting
				mock.AnythingOfType("*entity.SongFilters")). // filter struct
				Return(returnMockSongUsecase, tt.mockSUGetListSong.returnError).
				Once()

			handlers.NewSongHandler(router, mockSongusecase, mockMusicInfoUsecase, logger)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", tt.reqURL, strings.NewReader(string(filterBodyJSON)))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.NotEmpty(t, string(expectedInJSON))
			assert.Equal(t, string(expectedInJSON), w.Body.String())

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

// Test for POST /song/create
func TestCreateSong(t *testing.T) {
	ginTest()

	type mockMIUReturnVal struct {
		res *entity.SongRequest
		err error
	}
	type mockSUReturnalVal struct {
		res *entity.Song
		err error
	}
	d := time.Date(2024, 10, 6, 13, 13, 10, 0, time.Now().UTC().Location())
	tables := []struct {
		nameTest  string
		songInput *entity.SongRequest
		musicInfo *mockMIUReturnVal
		songUS    *mockSUReturnalVal
	}{
		{
			nameTest: "success",
			songInput: &entity.SongRequest{
				Song:  "Test song",
				Group: "Test group",
			},
			musicInfo: &mockMIUReturnVal{
				res: &entity.SongRequest{
					Song:        "Test song",
					Group:       "Test group",
					Link:        "Test link",
					Text:        "Test text",
					ReleaseDate: &d,
				},
				err: nil,
			},
			songUS: &mockSUReturnalVal{
				res: &entity.Song{
					ID:          1,
					Name:        "Test song",
					Link:        "Test group",
					Text:        "Test link",
					ReleaseDate: d,
					GroupID:     int64(1),
				},
				err: nil,
			},
		},
	}

	for _, tt := range tables {
		t.Run(tt.nameTest, func(t *testing.T) {
			router := gin.Default()
			mockSongusecase := new(mocks.SongUsecase)
			mockMusicInfoUsecase := new(mocks.MusicInfoUsecase)
			expectedResult := tt.songUS.res
			songReqJSON, err := json.Marshal(tt.songInput)
			assert.NoError(t, err)

			mockMusicInfoUsecase.On("GetInfo", mock.Anything, tt.songInput).Run(func(args mock.Arguments) {
				songReqFN := args.Get(1).(*entity.SongRequest)
				songReqFN.Link = tt.musicInfo.res.Link
				songReqFN.Text = tt.musicInfo.res.Text
				songReqFN.ReleaseDate = tt.musicInfo.res.ReleaseDate
			}).Return(tt.musicInfo.err).Once()
			mockSongusecase.On("CreateNewSong", mock.Anything, tt.musicInfo.res).Return(tt.songUS.res, tt.songUS.err).Once()

			handlers.NewSongHandler(router, mockSongusecase, mockMusicInfoUsecase, logger)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/song/create", strings.NewReader(string(songReqJSON)))

			expectedInJSON, _ := json.Marshal(expectedResult)
			router.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
			assert.Equal(t, string(expectedInJSON), w.Body.String())

			result := &entity.Song{}
			err = json.Unmarshal(w.Body.Bytes(), result)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, expectedResult, result)
		})
	}
}
