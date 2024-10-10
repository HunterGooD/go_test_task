package interfaces

import (
	"context"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name SongRepository
type SongRepository interface {
	WithTransaction(tx *sqlx.Tx) SongRepository

	CreateSong(ctx context.Context, group_id int64, songInput *entity.SongRequest) (*entity.Song, error)

	GetListSong(ctx context.Context, offset, limit int, with_deleted bool, filters *entity.SongFilters) ([]entity.Song, error)
	GetReverseListSongs(ctx context.Context, offset, limit int, with_deleted bool, filters *entity.SongFilters) ([]entity.Song, error)
	GetSongTextByID(ctx context.Context, songID int64) (string, error)
	GetByNames(ctx context.Context, song_name, group_name string) (*entity.Song, error)
	GetByName(ctx context.Context, song_name string, group_id int64) (*entity.Song, error)
	GetByID(ctx context.Context, id int64) (*entity.Song, error)
	Total(ctx context.Context, with_deleted bool, filter *entity.SongFilters) (int, error)

	UpdateFromMapByID(ctx context.Context, id int64, song *entity.Song, fields map[string]any) error

	DeleteSoftByID(ctx context.Context, id int64) error
	DeleteSoftByGroupID(ctx context.Context, id int64) error
	DeleteForceByID(ctx context.Context, id int64) error
	DeleteForceByGroupID(ctx context.Context, id int64) error
}
