package repository

import (
	"context"
	"fmt"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
	"github.com/HunterGooD/go_test_task/internal/utils"
	"github.com/jmoiron/sqlx"
)

type songRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewSongRepository(conn *sqlx.DB) *songRepository {
	return &songRepository{conn, nil}
}

func (s *songRepository) WithTransaction(tx *sqlx.Tx) interfaces.SongRepository {
	return &songRepository{s.db, tx}
}

// GetListSong filters any fields with values for where sql
// map is expected value Name|Link|Text|ReleaseDate|GroupID
// filters ? Why not use map ?? XD
func (s *songRepository) GetListSong(ctx context.Context, offset, limit int, filters *entity.SongFilters) ([]entity.Song, error) {
	res := make([]entity.Song, 0)
	whereStatemant := ""
	if filters != nil {
		utils.GetFilterString(filters) //TODO: change package without heavy depends
	}
	query := fmt.Sprintf(`SELECT s.*, g.id AS "group.id", g.g_name AS "group.g_name",
		 	g.created_at AS "group.created_at", g.update_at AS "group.update_at",
			g.deleted_at AS "group.deleted_at" 
		FROM groups g Inner join songs s ON g.id=s.group_id 
		WHERE s.deleted_at IS NULL AND g.deleted_at IS NULL 
		%s
		LIMIT $1 OFFSET $2;`, whereStatemant)

	err := s.db.SelectContext(ctx, &res, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (s *songRepository) Total(ctx context.Context) (int, error) {
	var total int
	query := "SELECT COUNT(*) FROM songs"
	err := s.db.GetContext(ctx, &total, query)
	return total, err
}

func (s *songRepository) GetSongTextByID(ctx context.Context, songID int64) (string, error) {
	var text string
	query := `SELECT m_text FROM songs WHERE id = $1 AND deleted_at IS NULL`
	err := s.db.GetContext(ctx, &text, query, songID)
	if err != nil {
		return "", err
	}
	return text, err
}

func (s *songRepository) GetSongTextByName(ctx context.Context, name string) (*entity.Song, error) {

	return nil, nil
}

func (s *songRepository) GetByName(ctx context.Context, name string) (*entity.Song, error) {
	return nil, nil
}

func (s *songRepository) GetByID(ctx context.Context, id int64) (*entity.Song, error) {
	return nil, nil
}

func (s *songRepository) UpdateFromMap(ctx context.Context, fields map[string]string) (*entity.Song, error) {
	return nil, nil
}

func (s *songRepository) DeleteSoftByName(ctx context.Context, name string) error {
	return nil
}

func (s *songRepository) DeleteSoftByID(ctx context.Context, id int64) error {
	return nil
}

func (s *songRepository) DeleteSoftByGroupID(ctx context.Context, id int64) error {
	return nil
}

func (s *songRepository) DeleteSoftSong(ctx context.Context) error {
	_ = `-- SOFT DELETE SONG
	UPDATE songs
		SET deleted_at=CURRENT_TIMESTAMP
		WHERE id=? OR name='?';`
	return nil
}

func (s *songRepository) DeleteForceByName(ctx context.Context, name string) error {
	return nil
}

func (s *songRepository) DeleteForceByID(ctx context.Context, id int64) error {
	return nil
}
