package types

// KVStore is the base interface for all methods related to a store. See the store package
type KVStore interface {
	// Get from the cache or tree
	Get(key []byte) []byte
	// Set to the cache or tree
	Set(key, value []byte)
	// IterateKeyRange over the tree
	IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool
}

// Cache extends KVStore adding an additional method to implement on a cache
type Cache interface {
	KVStore
	// Dump the cache to the tree
	ApplyToState()
}

// Store extends KVStore
type Store interface {
	KVStore
	// Commit is called on Abci commit to commit the tree to storage and update the hash
	Commit() CommitInfo
	// Close the store
	Close()
	// Refresh the check/deliver caches
	RefreshCache() Cache
}
