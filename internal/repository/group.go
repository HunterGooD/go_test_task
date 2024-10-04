package repository

import "github.com/jmoiron/sqlx"

type GroupRepository struct {
	Conn *sqlx.DB
}

func NewGroupRepository(conn *sqlx.DB) *GroupRepository {
	return &GroupRepository{conn}
}
