package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/HunterGooD/go_test_task/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetById1(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "g_name", "created_at", "update_at", "deleted_at"}).
		AddRow(1, "friks", time.Now(), time.Now(), nil)

	// query := "SELECT id, g_name, created_at, update_at, deleted_at FROM groups WHERE id=\\?"
	query := "SELECT \\* FROM groups WHERE id = \\$1"
	// prep := mock.ExpectPrepare(query)
	groupID := int64(1)
	mock.ExpectQuery(query).WithArgs(groupID).WillReturnRows(rows)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	group := repository.NewGroupRepository(sqlxDB)

	groupStruct, err := group.GetByID(context.TODO(), groupID)

	assert.NoError(t, err)
	assert.NotNil(t, groupStruct)
	assert.Equal(t, int64(1), groupStruct.ID)
	assert.Equal(t, "friks", groupStruct.GName)
	assert.Nil(t, groupStruct.DeletedAt)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetById2(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "g_name", "created_at", "update_at", "deleted_at"}).
		AddRow(1, "friks", time.Now(), time.Now(), time.Now())

	// query := "SELECT id, g_name, created_at, update_at, deleted_at FROM groups WHERE id=\\?"
	query := "SELECT \\* FROM groups WHERE id = \\$1"
	// prep := mock.ExpectPrepare(query)
	groupID := int64(1)
	mock.ExpectQuery(query).WithArgs(groupID).WillReturnRows(rows)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	group := repository.NewGroupRepository(sqlxDB)

	groupStruct, err := group.GetByID(context.TODO(), groupID)

	assert.NoError(t, err)
	assert.NotNil(t, groupStruct)
	assert.Equal(t, int64(1), groupStruct.ID)
	assert.Equal(t, "friks", groupStruct.GName)
	assert.NotNil(t, groupStruct.DeletedAt)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
