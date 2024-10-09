package repository

import (
	"context"
	"database/sql"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
	"github.com/jmoiron/sqlx"
)

type groupRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewGroupRepository(conn *sqlx.DB) *groupRepository {
	return &groupRepository{conn, nil}
}

func (g *groupRepository) WithTransaction(tx *sqlx.Tx) interfaces.GroupRepository {
	return &groupRepository{g.db, tx}
}

func (g *groupRepository) CreateGroup(ctx context.Context, group_name string) (*entity.Group, error) {
	var err error
	songReturn := &entity.Group{}

	query := `INSERT INTO groups(g_name )
		VALUES ($1) RETURNING id, g_name`
	// if transaction activ exec in transaction else db exec
	if g.tx != nil {
		err = g.tx.GetContext(ctx, songReturn, query, group_name)
	} else {
		err = g.db.GetContext(ctx, songReturn, query, group_name)
	}
	return songReturn, err
}

func (g *groupRepository) GetByName(ctx context.Context, group_name string) (*entity.Group, error) {
	groupResult := &entity.Group{}
	var err error
	query := `SELECT * FROM groups WHERE g_name = $1`

	if g.tx != nil {
		err = g.tx.GetContext(ctx, groupResult, query, group_name)
	} else {
		err = g.db.GetContext(ctx, groupResult, query, group_name)
	}
	if err == sql.ErrNoRows {
		err = entity.ErrNotFound
	}
	return groupResult, err
}

// GetByID if withDeleted is true then view all rows
func (g *groupRepository) GetByID(ctx context.Context, id int64) (*entity.Group, error) {
	res := &entity.Group{}
	sql := `SELECT * FROM groups WHERE id = $1`

	err := g.db.Get(res, sql, id)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (g *groupRepository) UpdateFromMap(ctx context.Context, fields map[string]string) (*entity.Group, error) {
	return nil, nil
}

func (g *groupRepository) DeleteSoftByName(ctx context.Context, name string) error {
	return nil
}

func (g *groupRepository) DeleteSoftByID(ctx context.Context, id int64) error {
	return nil
}

func (g *groupRepository) DeleteForceByName(ctx context.Context, name string) error {
	return nil
}

func (g *groupRepository) DeleteForceByID(ctx context.Context, id int64) error {
	return nil
}
