package types

import fmt "fmt"

// KVStore is the base interface for all methods related to a store. See the store package
type KVStore interface {
	// Get from the cache or tree
	Get(key []byte) ([]byte, error)
	// Get only from committed data
	GetCommitted(key []byte) ([]byte, error)
	// Set to the cache or tree
	Set(key, value []byte)
	// Delete a key/value pair
	Delete(key []byte)
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

// Creates the store prefix key.  Storage is scoped by the
// service name
func createServiceKey(serviceName string, key []byte) []byte {
	return append([]byte(fmt.Sprintf("%s::", serviceName)), key...)
}

// Store for Service Txs
type RWStore struct {
	store       Cache
	serviceName string
}

func NewRWStore(service string, store Cache) RWStore {
	return RWStore{
		store:       store,
		serviceName: service,
	}
}

func (rw RWStore) Insert(key, value []byte) {
	prefixedKey := createServiceKey(rw.serviceName, key)
	rw.store.Set(prefixedKey, value)
}

func (rw RWStore) Get(key []byte) ([]byte, error) {
	prefixedKey := createServiceKey(rw.serviceName, key)
	return rw.store.Get(prefixedKey)
}

func (rw RWStore) Delete(key []byte) {
	prefixedKey := createServiceKey(rw.serviceName, key)
	rw.store.Delete(prefixedKey)
}

type QueryStore struct {
	store       Cache
	serviceName string
}

func NewQueryStore(service string, store Cache) QueryStore {
	return QueryStore{
		store:       store,
		serviceName: service,
	}
}

func (qs QueryStore) Get(key []byte) ([]byte, error) {
	prefixedKey := createServiceKey(qs.serviceName, key)
	return qs.store.GetCommitted(prefixedKey)
}
