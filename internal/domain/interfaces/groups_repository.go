package interfaces

import (
	"context"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name GroupRepository
type GroupRepository interface {
	WithTransaction(tx *sqlx.Tx) GroupRepository

	CreateGroup(ctx context.Context, group_name string) (*entity.Group, error)

	GetByName(ctx context.Context, name string) (*entity.Group, error)
	GetByID(ctx context.Context, id int64) (*entity.Group, error)

	UpdateFromMap(ctx context.Context, fields map[string]string) (*entity.Group, error)

	DeleteSoftByName(ctx context.Context, name string) error
	DeleteSoftByID(ctx context.Context, id int64) error
	DeleteForceByName(ctx context.Context, name string) error
	DeleteForceByID(ctx context.Context, id int64) error
}
