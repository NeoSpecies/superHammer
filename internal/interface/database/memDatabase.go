package database

type IDatabase interface {
	LoadData() error
	Get(key string) (string, bool)
	Put(key string, value string) error
	Delete(key string)
	Persist() error
	BeginTransaction() error
	CommitTransaction() error
	RollbackTransaction()
}
