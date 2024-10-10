package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HunterGooD/go_test_task/config"
	"github.com/HunterGooD/go_test_task/internal/repository"
	"github.com/HunterGooD/go_test_task/internal/repository/transaction"
	"github.com/HunterGooD/go_test_task/internal/rest/handlers"
	"github.com/HunterGooD/go_test_task/internal/rest/middleware"
	"github.com/HunterGooD/go_test_task/internal/usecase"
	"github.com/HunterGooD/go_test_task/pkg/api"
	"github.com/HunterGooD/go_test_task/pkg/utils/logger"
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
	// init config from env
	cfg := config.NewConfig()

	// init logger
	logerQ := logger.NewJsonSlogLogger(os.Stdout, "info")
	logerQ.Debug("", map[string]any{
		"qwe":    1,
		"qwe123": "asdasd",
	})
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

	// init database connection
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

	// init and add midleware for router api
	// gin.SetMode(gin.DebugMode) TODO: set mode release and debug change
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

	// service for music info requests
	var musicInfoUsecase handlers.MusicInfoUsecase
	// if not set addres for server using mock impl
	if len(cfg.AddrMusicInfoService) != 0 {
		apiMusicInfo, err := api.NewClient(cfg.AddrMusicInfoService)
		if err != nil {
			slog.Error("error init client", slog.String("error on create client", err.Error()))
			panic(err)
		}
		musicInfoUsecase = usecase.NewMusicInfoUsecase(apiMusicInfo)
	} else {
		musicInfoUsecase = usecase.NewMusicInfoUsecaseImited()
	}

	// init handlers with usecase
	handlers.NewSongHandler(r, songUsecase, musicInfoUsecase)
	handlers.NewGroupHandler(r, groupUsecase)

	// init swagger docs
	handlers.NewSwaggerHandler(r)

	// run serving
	ctxApp, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen error", slog.Any("err", err))
			panic(err)
		}
		cancelApp()
	}()

	// gracefull stop
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// wait sys call
	<-quit

	slog.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(ctxApp, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", slog.Any("err", err))
		panic(err)
	}

	// catching ctx.Done(). timeout of 5 seconds.

	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		slog.Info("timeout of 5 seconds.")
	}

	slog.Info("Server exiting")
}
