package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/HunterGooD/go_test_task/config"
	"github.com/HunterGooD/go_test_task/internal/repository"
	"github.com/HunterGooD/go_test_task/internal/repository/transaction"
	"github.com/HunterGooD/go_test_task/internal/rest/handlers"
	"github.com/HunterGooD/go_test_task/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.NewConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
	r := gin.Default()
	// TODO: middleware

	// init repository for db execution
	songRepo := repository.NewSongRepository(dbconn)
	groupRepo := repository.NewGroupRepository(dbconn)

	// init transaction manager for transaction control with repository
	txManagerSongsGroups := transaction.NewTransactionManagerSongsGroups(dbconn, groupRepo, songRepo)

	// init usecases
	songUsecase := usecase.NewSongUsecase(songRepo, txManagerSongsGroups)
	groupUsecase := usecase.NewGroupUsecase(groupRepo, txManagerSongsGroups)

	// init handlers with usecase
	handlers.NewSongHandler(r, songUsecase, logger)
	handlers.NewGroupHandler(r, groupUsecase, logger)

	// run serving
	r.Run(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
}
