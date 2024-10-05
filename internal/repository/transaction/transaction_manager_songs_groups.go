package transaction

import (
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
	"github.com/jmoiron/sqlx"
)

type TransactionManagerSongsGroups struct {
	db *sqlx.DB
	tx *sqlx.Tx

	group interfaces.GroupRepository
	song  interfaces.SongRepository
}

func NewTransactionManagerSongsGroups(db *sqlx.DB, group interfaces.GroupRepository, song interfaces.SongRepository) *TransactionManagerSongsGroups {
	return &TransactionManagerSongsGroups{db, nil, group, song}
}

func (t *TransactionManagerSongsGroups) SongRepository() interfaces.SongRepository {
	return t.song.WithTransaction(t.tx)
}

func (t *TransactionManagerSongsGroups) GroupRepository() interfaces.GroupRepository {
	return t.group.WithTransaction(t.tx)
}

func (t *TransactionManagerSongsGroups) Begin() error {
	var err error
	t.tx, err = t.db.Beginx()
	return err
}

func (t *TransactionManagerSongsGroups) Commit() error {
	return t.tx.Commit()
}

func (t *TransactionManagerSongsGroups) Rollback() error {
	return t.tx.Rollback()
}
