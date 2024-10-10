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

	log interfaces.Logger
}

func NewGroupRepository(conn *sqlx.DB, logger interfaces.Logger) *groupRepository {
	logger.Info("Creat group repo")
	return &groupRepository{conn, nil, logger}
}

func (g *groupRepository) WithTransaction(tx *sqlx.Tx) interfaces.GroupRepository {
	g.log.Info("Creat group repo with transaction")
	return &groupRepository{g.db, tx, g.log}
}

func (g *groupRepository) CreateGroup(ctx context.Context, group_name string) (*entity.Group, error) {
	var err error
	group := &entity.Group{}
	query := `INSERT INTO groups(g_name )
		VALUES ($1) RETURNING id, g_name`

	// if transaction activ exec in transaction else db exec
	if g.tx != nil {
		g.log.Info("create group with transaction", map[string]any{
			"group_name": group_name,
			"query":      query,
		})
		err = g.tx.GetContext(ctx, group, query, group_name)
	} else {
		g.log.Info("create group without transaction", map[string]any{
			"group_name": group_name,
			"query":      query,
		})
		err = g.db.GetContext(ctx, group, query, group_name)
	}

	if err != nil {
		g.log.Error("group creating error", map[string]any{
			"err": err,
		})
		return nil, err
	}

	g.log.Info("create group success", map[string]any{
		"group_name": group_name,
		"group_id":   group.ID,
	})
	g.log.Debug("create group success", map[string]any{
		"group": group,
	})

	return group, nil
}

func (g *groupRepository) GetByName(ctx context.Context, group_name string) (*entity.Group, error) {
	group := &entity.Group{}
	var err error
	query := `SELECT * FROM groups WHERE g_name = $1`

	if g.tx != nil {
		g.log.Info("get group with transaction", map[string]any{
			"group_name": group_name,
			"query":      query,
		})
		err = g.tx.GetContext(ctx, group, query, group_name)
	} else {
		g.log.Info("create group without transaction", map[string]any{
			"group_name": group_name,
			"query":      query,
		})
		err = g.db.GetContext(ctx, group, query, group_name)
	}

	if err != nil {
		g.log.Error("group get error", map[string]any{
			"err": err,
		})
		if err == sql.ErrNoRows {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	g.log.Info("get group success", map[string]any{
		"group_name": group_name,
		"group_id":   group.ID,
	})
	g.log.Debug("get group success", map[string]any{
		"group": group,
	})

	return group, nil
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
