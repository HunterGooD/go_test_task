package interfaces

//go:generate mockery --name TransactionManagerSongsGroups
type TransactionManagerSongsGroups interface {
	Begin() error
	Commit() error
	Rollback() error
	SongRepository() SongRepository
	GroupRepository() GroupRepository
}
