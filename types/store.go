package types

import (
	"errors"

	"github.com/davebryson/menta/storage"
	"github.com/davebryson/menta/store"
)

// KVStore re-exported for the sdk
type (
	KVStore  = store.KVStore
	Snapshot = storage.TreeReader
	Cache    = storage.Cache
)

func PrefixedKey(service, key []byte) []byte {
	res := make([]byte, len(service)+len(key))
	copy(res, service)
	copy(res[len(service):], key)
	return res
}

type PrefixedKVStore struct {
	prefix []byte
	store  Cache
}

func NewPrefixedKVStore(prefix string, store Cache) PrefixedKVStore {
	return PrefixedKVStore{
		prefix: []byte(prefix),
		store:  store,
	}
}

func (ps PrefixedKVStore) key(k []byte) []byte {
	return PrefixedKey(ps.prefix, k)
}

func (ps PrefixedKVStore) Get(key []byte) ([]byte, error) {
	data, err := ps.store.Get(ps.key(key))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ps PrefixedKVStore) Has(key []byte) bool {
	_, err := ps.Get(key)
	return err == nil
}

func (ps PrefixedKVStore) Put(key, value []byte) error {
	if key == nil {
		return errors.New("PrefixKVStore: Key is nil")
	}
	if value == nil {
		return errors.New("PrefixKVStore: Value is nil")
	}
	ps.store.Put(ps.key(key), value)
	return nil
}

func (ps PrefixedKVStore) Remove(key []byte) {
	ps.store.Remove(ps.key(key))
}

type PrefixedSnapshot struct {
	prefix []byte
	store  Snapshot
}

func NewPrefixedSnapshot(prefix string, snapshot Snapshot) PrefixedSnapshot {
	return PrefixedSnapshot{
		prefix: []byte(prefix),
		store:  snapshot,
	}
}

func (ps PrefixedSnapshot) Get(key []byte) ([]byte, error) {
	return ps.store.Get(PrefixedKey(ps.prefix, key))
}
