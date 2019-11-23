package types

type StoreReader interface {
	// Get from the cache or tree
	Get(key []byte) ([]byte, error)
	// Iterateover the tree/db
	Iterate(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool
}

type StoreWriter interface {
	// Set to the cache or tree
	Set(key, value []byte)
	// Delete a key/value
	Delete(key []byte)
}

// KVStore is the base interface for all methods related to a store. See the store package
type RWStore interface {
	StoreReader
	StoreWriter
}

// Cache extends KVStore adding an additional method to implement on a cache
type Cache interface {
	RWStore
	// Dump the cache to the tree
	ApplyToState()
}

// Store extends KVStore
type Store interface {
	RWStore
	// Commit is called on Abci commit to commit the tree to storage and update the hash
	Commit() CommitInfo
	// Close the store
	Close()
	// Refresh the check/deliver caches
	RefreshCache() Cache
}
