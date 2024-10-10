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
	log   interfaces.Logger
}

func NewTransactionManagerSongsGroups(db *sqlx.DB, group interfaces.GroupRepository, song interfaces.SongRepository, logger interfaces.Logger) *TransactionManagerSongsGroups {
	logger.Info("Create transaction manager")
	return &TransactionManagerSongsGroups{db, nil, group, song, logger}
}

func (t *TransactionManagerSongsGroups) SongRepository() interfaces.SongRepository {
	t.log.Info("Get SongRepository with transaction")
	return t.song.WithTransaction(t.tx)
}

func (t *TransactionManagerSongsGroups) GroupRepository() interfaces.GroupRepository {
	t.log.Info("Get GroupRepository with transaction")
	return t.group.WithTransaction(t.tx)
}

func (t *TransactionManagerSongsGroups) Begin() error {
	var err error

	t.log.Info("Begin transaction")

	t.tx, err = t.db.Beginx()
	if err != nil {
		t.log.Error("Error on begin transaction", map[string]any{
			"err": err,
		})
		return err
	}

	return nil
}

func (t *TransactionManagerSongsGroups) Commit() error {
	t.log.Info("Commit transaction")
	return t.tx.Commit()
}

func (t *TransactionManagerSongsGroups) Rollback() error {
	t.log.Warn("Rollback transaction")
	return t.tx.Rollback()
}
