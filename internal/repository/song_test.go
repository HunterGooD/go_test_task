package repository_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/HunterGooD/go_test_task/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func TestGetListSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	songRepo := repository.NewSongRepository(sqlxDB, logger)

	rows := sqlmock.NewRows([]string{"id", "m_name", "m_link", "m_text", "m_release_date", "group_id", "created_at", "update_at", "deleted_at", "group.id", "group.g_name", "group.created_at", "group.update_at", "group.deleted_at"}).
		AddRow(1, "name song 1", "hello", "qweasdzxc123", time.Now(), 1, time.Now(), time.Now(), nil, 1, "friks", time.Now(), time.Now(), nil).
		AddRow(2, "name song 2", "world", "qweasdzxc123", time.Now(), 1, time.Now(), time.Now(), nil, 1, "friks", time.Now(), time.Now(), nil)
	query := `SELECT 
			s\.\*,
			g\.id AS \"group\.id\",
			g\.g_name AS \"group\.g_name\",
			g\.created_at AS \"group\.created_at\",
			g\.update_at AS \"group\.update_at\",
			g\.deleted_at AS \"group\.deleted_at\" 
		FROM groups g Inner join songs s ON g\.id=s\.group_id 
		WHERE s\.deleted_at IS NULL AND g\.deleted_at IS NULL 
		LIMIT \$1 OFFSET \$2;`

	limit, offset := 10, 0
	// filters := make(map[string]string)
	// filters["m_name"] = "name song 1"
	// filters["m_link"] = "hello"

	mock.ExpectQuery(query).
		WithArgs(limit, offset).
		WillReturnRows(rows)

	songsList, err := songRepo.GetListSong(context.TODO(), offset, limit, false, nil)

	assert.NoError(t, err)
	assert.NotNil(t, songsList)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSongTextByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	songRepo := repository.NewSongRepository(sqlxDB, logger)

	text := "asdkasndkjasjkdhajksdakjjsdkjasjkdnajksndajksndnjkansjkdbajksbd"
	// rows := sqlmock.NewRows([]string{"id", "m_name", "m_link", "m_text", "m_release_date", "created_at", "update_at", "deleted_at"}).
	rows := sqlmock.NewRows([]string{"m_text"}).
		AddRow(text)
	query := `SELECT m_text FROM songs WHERE id = \$1 AND deleted_at IS NULL`
	songID := int64(1)
	mock.ExpectQuery(query).
		WithArgs(songID).
		WillReturnRows(rows)

	songRes, err := songRepo.GetSongTextByID(context.TODO(), songID)
	assert.NoError(t, err)
	assert.NotEmpty(t, songRes)
	assert.Equal(t, "asdkasndkjasjkdhajksdakjjsdkjasjkdnajksndajksndnjkansjkdbajksbd", songRes)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSongTextByIDError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	songRepo := repository.NewSongRepository(sqlxDB, logger)

	// rows := sqlmock.NewRows([]string{"id", "m_name", "m_link", "m_text", "m_release_date", "created_at", "update_at", "deleted_at"}).

	query := `SELECT m_text FROM songs WHERE id = \$1 AND deleted_at IS NULL`
	songID := int64(1)
	mock.ExpectQuery(query).
		WithArgs(songID).
		WillReturnError(sql.ErrNoRows)

	songRes, err := songRepo.GetSongTextByID(context.TODO(), songID)
	assert.Empty(t, songRes)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
