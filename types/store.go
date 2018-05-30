package types

type KVStore interface {
	Get(key []byte) []byte
	Set(key, value []byte)
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
