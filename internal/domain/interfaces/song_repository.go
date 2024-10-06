package interfaces

import (
	"context"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name SongRepository
type SongRepository interface {
	WithTransaction(tx *sqlx.Tx) SongRepository

	GetListSong(ctx context.Context, offset, limit int, filters *entity.SongFilters) ([]entity.Song, error)
	GetSongTextByID(ctx context.Context, songID int64) (string, error)
	GetSongTextByName(ctx context.Context, name string) (*entity.Song, error)
	GetByName(ctx context.Context, name string) (*entity.Song, error)
	GetByID(ctx context.Context, id int64) (*entity.Song, error)
	Total(ctx context.Context) (int, error)
	UpdateFromMap(ctx context.Context, fields map[string]string) (*entity.Song, error)
	DeleteSoftByName(ctx context.Context, name string) error
	DeleteSoftByID(ctx context.Context, id int64) error
	DeleteSoftByGroupID(ctx context.Context, id int64) error
	DeleteSoftSong(ctx context.Context) error
	DeleteForceByName(ctx context.Context, name string) error
	DeleteForceByID(ctx context.Context, id int64) error
}
