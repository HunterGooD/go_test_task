package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/HunterGooD/go_test_task/config"
	"github.com/HunterGooD/go_test_task/internal/repository"
	"github.com/HunterGooD/go_test_task/internal/repository/transaction"
	"github.com/HunterGooD/go_test_task/internal/rest/handlers"
	"github.com/HunterGooD/go_test_task/internal/rest/middleware"
	"github.com/HunterGooD/go_test_task/internal/usecase"
	"github.com/HunterGooD/go_test_task/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// @title Music librarry API
// @version 0.0.1alpha
// @description This is a simple api for music librarry
// @termsOfService http://swagger.io/terms/

// @BasePath /
func main() {
	cfg := config.NewConfig()
	var opts *slog.HandlerOptions
	switch cfg.LogLevel {
	case "INFO":
		opts = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	case "DEBUG":
		opts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	slog.SetDefault(logger)

	// TODO: switch cfg.DBtype for any connection db, add drivers for mysql, sqlite3
	dbconn, err := sqlx.Open(cfg.DBType, cfg.DSN)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	if err := dbconn.Ping(); err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	// router for api
	r := gin.New()
	// TODO: middleware
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// init repository for db execution
	songRepo := repository.NewSongRepository(dbconn)
	groupRepo := repository.NewGroupRepository(dbconn)

	// init transaction manager for transaction control with repository
	txManagerSongsGroups := transaction.NewTransactionManagerSongsGroups(dbconn, groupRepo, songRepo)

	// init usecases
	songUsecase := usecase.NewSongUsecase(songRepo, txManagerSongsGroups)
	groupUsecase := usecase.NewGroupUsecase(groupRepo, txManagerSongsGroups)
	apiMusicInfo, _ := api.NewClient("")
	musicInfoUsecase := usecase.NewMusicInfoUsecase(apiMusicInfo)

	// init handlers with usecase
	handlers.NewSongHandler(r, songUsecase, musicInfoUsecase)
	handlers.NewGroupHandler(r, groupUsecase)

	// init swagger docs
	handlers.NewSwaggerHandler(r)
	// run serving
	r.Run(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
}
