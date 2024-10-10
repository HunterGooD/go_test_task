package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	slogLogger := logger.NewJsonSlogLogger(os.Stdout, strings.ToUpper(cfg.LogLevel))
	slogLogger.Info("init logger, config success!")
	slogLogger.Debug("init config", map[string]any{"config": cfg})

	// init database connection
	// TODO: switch cfg.DBtype for any connection db, add drivers for mysql, sqlite3
	slogLogger.Info("initial database connection ...", map[string]any{
		"database_type": cfg.DBType,
		"dsn":           cfg.DSN,
	})
	dbconn, err := sqlx.Open(cfg.DBType, cfg.DSN)
	if err != nil {
		slogLogger.Error(err.Error())
		panic(err)
	}

	if err := dbconn.Ping(); err != nil {
		slogLogger.Error(err.Error())
		panic(err)
	}
	slogLogger.Info("database connection success")

	// init and add midleware for router api
	// gin.SetMode(gin.DebugMode) TODO: set mode release and debug change
	r := gin.New()

	// TODO: middleware
	slogLogger.Info("gin middleware init ...")

	midlWare := middleware.NewMiddleware(slogLogger)
	r.Use(midlWare.Logger())
	r.Use(gin.Recovery())

	slogLogger.Info("gin middleware success")

	// init repository for db execution
	slogLogger.Info("repository init ...")

	songRepo := repository.NewSongRepository(dbconn, slogLogger)
	groupRepo := repository.NewGroupRepository(dbconn, slogLogger)

	// init transaction manager for transaction control with repository
	txManagerSongsGroups := transaction.NewTransactionManagerSongsGroups(dbconn, groupRepo, songRepo, slogLogger)

	slogLogger.Info("repository init success")

	// init usecases
	slogLogger.Info("usecases init ...")

	songUsecase := usecase.NewSongUsecase(songRepo, txManagerSongsGroups, slogLogger)
	groupUsecase := usecase.NewGroupUsecase(groupRepo, txManagerSongsGroups, slogLogger)

	// service for music info requests
	var musicInfoUsecase handlers.MusicInfoUsecase
	// if not set addres for server using mock impl
	if len(cfg.AddrMusicInfoService) != 0 {
		apiMusicInfo, err := api.NewClient(cfg.AddrMusicInfoService)
		if err != nil {
			slogLogger.Error("error init client", map[string]any{
				"error on create client": err.Error(),
			})
			panic(err)
		}
		musicInfoUsecase = usecase.NewMusicInfoUsecase(apiMusicInfo, slogLogger)
	} else {
		musicInfoUsecase = usecase.NewMusicInfoUsecaseImited(slogLogger)
	}
	slogLogger.Info("usecases init success")

	// init handlers with usecase
	slogLogger.Info("router register paths ...")

	handlers.NewSongHandler(r, songUsecase, musicInfoUsecase, slogLogger)
	handlers.NewGroupHandler(r, groupUsecase, slogLogger)

	// init swagger docs
	handlers.NewSwaggerHandler(r, slogLogger)

	slogLogger.Info("router register paths success")

	// run serving
	slogLogger.Info("gin server starting ...")

	ctxApp, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: r.Handler(),
	}

	go func() {
		slogLogger.Info("gin server start")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slogLogger.Error("listen error", map[string]any{"err": err})
			panic(err)
		}
		cancelApp()
	}()

	// gracefull stop
	slogLogger.Info("init gracefull ...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	slogLogger.Info("gracefull done, wait signals", map[string]any{
		"syscall.SIGINT":  syscall.SIGINT,
		"syscall.SIGTERM": syscall.SIGTERM,
	})

	// wait system call
	slogLogger.Info("getting syscal", map[string]any{"syscal": <-quit})

	slogLogger.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(ctxApp, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slogLogger.Error("Server Shutdown:", map[string]any{"err": err})
		panic(err)
	}

	// catching ctx.Done(). timeout of 5 seconds.

	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		slogLogger.Info("timeout of 5 seconds.")
	}

	slogLogger.Info("Server exiting") // BUG: service stopping with status 1 and dong exit with 0
}
