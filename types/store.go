package types

import (
	"errors"

	"github.com/davebryson/menta/store"
)

// KVStore re-exported for the sdk
type KVStore = store.KVStore

// NamedStore is used to provide scoped, read/write access to storage.
// Keys in the store are automatically prefixed with the given prefix name.
type NamedStore struct {
	prefix []byte
	store  KVStore
}

// NewNamedStore creates an instance
func NewNamedStore(serviceName string, store KVStore) NamedStore {
	return NamedStore{
		prefix: []byte(serviceName),
		store:  store,
	}
}

func (ps NamedStore) key(key []byte) []byte {
	return prefixKey(ps.prefix, key)
}

// Get something from the store by key
func (ps NamedStore) Get(key []byte) ([]byte, error) {
	data, err := ps.store.Get(ps.key(key))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Query a value by key from committed state
func (ps NamedStore) Query(key []byte) ([]byte, error) {
	data, err := ps.store.GetCommitted(ps.key(key))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Has : does the store contain the given key?
func (ps NamedStore) Has(key []byte) bool {
	_, err := ps.Get(key)
	return err == nil
}

// Put a key/value into the store
func (ps NamedStore) Put(key, value []byte) error {
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
func (ps NamedStore) Delete(key []byte) {
	ps.store.Delete(ps.key(key))
}

// generate a prefixed key
func prefixKey(service, key []byte) []byte {
	res := make([]byte, len(service)+len(key))
	copy(res, service)
	copy(res[len(service):], key)
	return res
}
