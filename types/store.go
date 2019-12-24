package types

import (
	"errors"

	"github.com/davebryson/menta/store"
)

// re-export store types for 'sdk' usage
type (
	KVStore = store.KVStore
	Store   = store.Store
	Cache   = store.Cache
)

// PrefixedRWStore is used to provide scoped, read/write access to storage.
// Keys in the store are automatically prefixed with the given prefix name.
// For example, this is used by Services to ensure all storage keys are
// prefixed with the service name
type PrefixedRWStore struct {
	prefix []byte
	store  KVStore
}

// NewPrefixRW creates an instance
func NewPrefixRW(serviceName string, store KVStore) PrefixedRWStore {
	return PrefixedRWStore{
		prefix: []byte(serviceName),
		store:  store,
	}
}

func (ps PrefixedRWStore) key(key []byte) []byte {
	return prefixKey(ps.prefix, key)
}

// Get something from the store by key
func (ps PrefixedRWStore) Get(key []byte) ([]byte, error) {
	data, err := ps.store.Get(ps.key(key))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Has : does the store contain the given key?
func (ps PrefixedRWStore) Has(key []byte) bool {
	_, err := ps.Get(key)
	return err == nil
}

// Insert a key/value into the store
func (ps PrefixedRWStore) Insert(key, value []byte) error {
	if key == nil {
		return errors.New("PrefixStore: Key is nil")
	}
	if value == nil {
		return errors.New("PrefixStore: Value is nil")
	}
	ps.store.Set(ps.key(key), value)
	return nil
}

// Delete a key/value from the store
func (ps PrefixedRWStore) Delete(key []byte) {
	ps.store.Delete(ps.key(key))
}

// generate a prefixed key
func prefixKey(service, key []byte) []byte {
	res := make([]byte, len(service)+len(key))
	copy(res, service)
	copy(res[len(service):], key)
	return res
}

// PrefixedReadOnlyStore provided scoped, read-only access to storage.
// Data returned is from committed data (not the cache).
type PrefixedReadOnlyStore struct {
	prefix []byte
	store  KVStore
}

func (pr PrefixedReadOnlyStore) key(key []byte) []byte {
	return prefixKey(pr.prefix, key)
}

// NewPrefixReadOnly creates a new instance
func NewPrefixReadOnly(serviceName string, store KVStore) PrefixedReadOnlyStore {
	return PrefixedReadOnlyStore{
		prefix: []byte(serviceName),
		store:  store,
	}
}

// Get a value for the given key from committed data
func (pr PrefixedReadOnlyStore) Get(key []byte) ([]byte, error) {
	return pr.store.GetCommitted(pr.key(key))
}
