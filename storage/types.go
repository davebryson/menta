package storage

import "github.com/cosmos/iavl"

// TreeReader provides read access to committed state
type TreeReader interface {
	// Get from committed state in the tree
	Get(key []byte) ([]byte, error)
	// GetWithProof returns the value with a Proof
	GetWithProof(key []byte) ([]byte, *iavl.RangeProof, error)
	// IterateKeyRange over committed state
	IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool
}

// TreeWriter writes to state and provides a snapshot of committed state
type TreeWriter interface {
	Commit(batch map[string]CacheOp) CommitData
	Snapshot() TreeReader
}

// Cache is the used to batch writes for commit
type Cache interface {
	// Get from the cache or tree
	Get(key []byte) ([]byte, error)
	// Does the store contain the given key
	Has(key []byte) bool
	// Put to the cache or tree
	Put(key, value []byte)
	// Delete a key/value pair
	Remove(key []byte)
	// ToBatch returns the cache storage
	ToBatch() map[string]CacheOp
}
