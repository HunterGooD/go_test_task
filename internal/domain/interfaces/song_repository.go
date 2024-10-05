package interfaces

import (
	"context"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name SongRepository
type SongRepository interface {
	WithTransaction(tx *sqlx.Tx) SongRepository

	GetListSong(ctx context.Context, page, pageSize int, filters map[string]string) ([]entity.Song, error)
	GetSongTextByID(ctx context.Context, songID int64) (*entity.Song, error)
	GetSongTextByName(ctx context.Context, name string) (*entity.Song, error)
	GetByName(ctx context.Context, name string) (*entity.Song, error)
	GetByID(ctx context.Context, id int64) (*entity.Song, error)
	UpdateFromMap(ctx context.Context, fields map[string]string) (*entity.Song, error)
	DeleteSoftByName(ctx context.Context, name string) error
	DeleteSoftByID(ctx context.Context, id int64) error
	DeleteSoftByGroupID(ctx context.Context, id int64) error
	DeleteSoftSong(ctx context.Context) string
	DeleteForceByName(ctx context.Context, name string) error
	DeleteForceByID(ctx context.Context, id int64) error
}
