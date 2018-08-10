package types

// TODO: Add IterateKeyRange to this interface
// Implement on both Cache and Store
type KVStore interface {
	Get(key []byte) []byte
	Set(key, value []byte)
	IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool
}

type Cache interface {
	KVStore
	GetAccount(address []byte) (*Account, error)
	SetAccount(account *Account)
	ApplyToState()
}

type Store interface {
	KVStore
	Commit() CommitInfo
	Close()
	RefreshCache() Cache
}
