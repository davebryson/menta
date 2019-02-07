package store

// In-memory cache used by the menta app
import (
	"sort"

	"github.com/davebryson/menta/types"
)

// Implements KVCache
var _ types.Cache = (*KVCache)(nil)

// Used to track whether a k/v pair has been updated.
type dataCacheObj struct {
	value []byte
	dirty bool
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
	cache.storage[string(key)] = dataCacheObj{val, true}
}

// Get a value for a given key.  Try the cache first and then the state db
func (cache *KVCache) Get(key []byte) []byte {
	cacheKey := string(key)

	// check the cache
	data := cache.storage[cacheKey]
	if data.value != nil {
		return data.value
	}

	// Not in the cache, go to cold storage
	value := cache.statedb.Get(key)
	if value != nil {
		cache.storage[cacheKey] = dataCacheObj{value, false}
		return value
	}

	return nil
}

// IterateKeyRange returns results that are processed via the callback func
func (cache *KVCache) IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool {
	return cache.statedb.IterateKeyRange(start, end, ascending, fn)
}

// ApplyToState is called during abci.commit().   It sorts all keys in the cache
// for determinism, then writes the set to the tree
func (cache *KVCache) ApplyToState() {
	storageKeys := make([]string, 0, len(cache.storage))
	for k := range cache.storage {
		storageKeys = append(storageKeys, k)
	}

	sort.Strings(storageKeys)

	for _, k := range storageKeys {
		data := cache.storage[k]
		// We don't (re)set unchanged data
		if data.value == nil || !data.dirty {
			continue
		}
		cache.statedb.Set([]byte(k), data.value)
	}
}
