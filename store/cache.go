package store

// In-memory cache used by the menta app
import (
	"sort"

	sdk "github.com/davebryson/menta/types"
)

// Implements KVCache
var _ sdk.Cache = (*KVCache)(nil)

// Used to track whether a k/v pair has been updated.
type dataCacheObj struct {
	value  []byte
	dirty  bool
	delete bool
}

// KVCache used by the app. It wraps a simple cache and access to the State Store
type KVCache struct {
	statedb *StateStore
	storage map[string]dataCacheObj
}

// NewCache return a fresh empty cache with ref to the State Store
func NewCache(store *StateStore) *KVCache {
	return &KVCache{
		statedb: store,
		storage: make(map[string]dataCacheObj),
	}
}

// Set a key in the cache
func (cache *KVCache) Set(key, val []byte) {
	cache.storage[string(key)] = dataCacheObj{val, true, false}
}

// Delete a key/value
func (cache *KVCache) Delete(key []byte) {
	if cache.Exists(key) {
		cache.storage[string(key)] = dataCacheObj{nil, false, true}
	}
}

// Exists - checks for a given key
func (cache *KVCache) Exists(key []byte) bool {
	_, err := cache.Get(key)
	if err == nil {
		return true
	}
	return false
}

// GetCommitted only returns committed data, nothing cached
func (cache *KVCache) GetCommitted(key []byte) ([]byte, error) {
	return cache.statedb.Get(key)
}

// Get a value for a given key.  Try the cache first and then the state db
func (cache *KVCache) Get(key []byte) ([]byte, error) {
	cacheKey := string(key)

	// check the cache
	data, ok := cache.storage[cacheKey]
	if ok && !data.delete {
		return data.value, nil
	}

	// Not in the cache, go to cold storage
	value, err := cache.statedb.Get(key)
	if err == nil {
		// cache it as not-dirty
		cache.storage[cacheKey] = dataCacheObj{value, false, false}
		return value, nil
	}

	return nil, ErrValueNotFound
}

// IterateKeyRange returns results that are processed via the callback func
func (cache *KVCache) IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool {
	return cache.statedb.IterateKeyRange(start, end, ascending, fn)
}

// ApplyToState is called during abci.commit().   It sorts all keys in the cache
// for determinism, then writes the set to the tree
func (cache *KVCache) ApplyToState() {
	storageKeys := make([]string, 0, len(cache.storage))

	for key := range cache.storage {
		storageKeys = append(storageKeys, key)
	}

	// Sort keys for determinism (required by IAVL)
	sort.Strings(storageKeys)

	for _, key := range storageKeys {
		data := cache.storage[key]

		// do delete and continue
		if data.delete {
			cache.statedb.Delete([]byte(key))
			continue
		}

		// Only insert dirty data. We don't re-insert unchanged, cached data
		if data.dirty {
			cache.statedb.Set([]byte(key), data.value)
		}
	}
}
