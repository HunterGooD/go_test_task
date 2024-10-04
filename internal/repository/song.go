package repository

import (
	"context"

	"github.com/HunterGooD/go_test_task/internal/entity"
	"github.com/jmoiron/sqlx"
)

type SongRepository struct {
	Conn *sqlx.DB
}

func NewSongRepository(conn *sqlx.DB) *SongRepository {
	return &SongRepository{conn}
}

func (s *SongRepository) GetSongs(ctx context.Context, filters map[string]string) ([]entity.Song, error) {

	return nil, nil
}
