package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (s *songRepository) CreateSong(ctx context.Context, group_id int64, songInput *entity.SongRequest) (*entity.Song, error) {
	var err error
	var songReturn *entity.Song

	query := `INSERT INTO public.songs(
		m_name, m_link, m_text, m_release_date, group_id)
		VALUES ($1, $2, $3, $4, $5); RETURNING id, m_name, m_link, m_text, m_release_date, group_id`
	// if transaction activ exec in transaction else db exec
	if s.tx != nil {
		err = s.tx.GetContext(ctx, songReturn, query, songInput.Song, songInput.Link, songInput.Text, songInput.ReleaseDate, group_id)
	} else {
		err = s.db.GetContext(ctx, songReturn, query, songInput.Song, songInput.Link, songInput.Text, songInput.ReleaseDate, group_id)
	}
	return songReturn, err
}

// GetListSong filters any fields with values for where sql
// map is expected value Name|Link|Text|ReleaseDate|GroupID
// filters ? Why not use map ?? XD
func (s *songRepository) GetListSong(ctx context.Context, offset, limit int, filters *entity.SongFilters) ([]entity.Song, error) {
	res := make([]entity.Song, 0)
	whereStatemant := ""
	var params []any
	params = append(params, limit, offset)
	if filters != nil {
		var arg []any
		whereStatemant, arg = utils.GetFilterString(len(params)+1, filters) //TODO: change package without heavy depends
		if len(arg) > 0 {
			params = append(params, arg...)
		}
	}
	query := fmt.Sprintf(`SELECT s.*, g.id AS "group.id", g.g_name AS "group.g_name",
		 	g.created_at AS "group.created_at", g.update_at AS "group.update_at",
			g.deleted_at AS "group.deleted_at" 
		FROM groups g Inner join songs s ON g.id=s.group_id 
		WHERE s.deleted_at IS NULL AND g.deleted_at IS NULL 
		%s
		LIMIT $1 OFFSET $2;`, whereStatemant)

	err := s.db.SelectContext(ctx, &res, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return res, err
}

func (s *songRepository) Total(ctx context.Context) (int, error) {
	var total int
	query := "SELECT COUNT(*) FROM songs"
	err := s.db.GetContext(ctx, &total, query)
	if err == sql.ErrNoRows {
		return total, entity.ErrNotFound
	}
	return total, err
}

func (s *songRepository) GetSongTextByID(ctx context.Context, songID int64) (string, error) {
	var text string
	query := `SELECT m_text FROM songs WHERE id = $1 AND deleted_at IS NULL`
	err := s.db.GetContext(ctx, &text, query, songID)
	if err == sql.ErrNoRows {
		return "", entity.ErrNotFound
	}
	return text, err
}

func (s *songRepository) GetByNames(ctx context.Context, song_name, group_name string) (*entity.Song, error) {
	return nil, nil
}

func (s *songRepository) GetByName(ctx context.Context, song_name string, group_id int64) (*entity.Song, error) {

	return nil, nil
}

func (s *songRepository) GetByID(ctx context.Context, id int64) (*entity.Song, error) {
	return nil, nil
}

func (s *songRepository) UpdateFromMapByID(ctx context.Context, id int64, fields map[string]string) (*entity.Song, error) {
	return nil, nil
}

func (s *songRepository) UpdateFromMapByNames(ctx context.Context, song_name, group_name string, fields map[string]string) (*entity.Song, error) {
	return nil, nil
}

func (s *songRepository) DeleteSoftByID(ctx context.Context, id int64) error {
	query := `UPDATE songs
		SET deleted_at=CURRENT_TIMESTAMP
		WHERE id= $1;`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		if rows == 0 {
			return entity.ErrNotFound
		}
		return err
	}
	return err
}

func (s *songRepository) DeleteSoftByGroupID(ctx context.Context, id int64) error {
	return nil
}

func (s *songRepository) DeleteSoftByNames(ctx context.Context, song_name, group_name string) error {
	return nil
}

func (s *songRepository) DeleteSoftSong(ctx context.Context) error {
	_ = `-- SOFT DELETE SONG
	UPDATE songs
		SET deleted_at=CURRENT_TIMESTAMP
		WHERE id=? OR name='?';`
	return nil
}

func (s *songRepository) DeleteForceByID(ctx context.Context, id int64) error {
	return nil
}

func (s *songRepository) DeleteForceByGroupID(ctx context.Context, id int64) error {
	return nil
}
