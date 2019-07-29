package store

// In-memory cache used by the menta app
import (
	"sort"
	"sync"

	"github.com/davebryson/menta/types"
)

// Implements KVCache
var _ types.Cache = (*KVCache)(nil)

// Used to track whether a k/v pair has been updated/deleted.
type dataCacheObj struct {
	value  []byte
	dirty  bool
	delete bool
}

// KVCache wraps a simple cache and access to the state store
type KVCache struct {
	mtx     sync.Mutex
	statedb *StateStore
	storage map[string]dataCacheObj
}

// NewCache return a fresh empty cache with a ref to the State Store
func NewCache(store *StateStore) *KVCache {
	return &KVCache{
		statedb: store,
		storage: make(map[string]dataCacheObj),
	}
}

// Set a key in the cache
func (cache *KVCache) Set(key, val []byte) {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()

	cache.storage[string(key)] = dataCacheObj{val, true, false}
}

// Delete a k/v
func (cache *KVCache) Delete(key []byte) {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()

	cacheKey := string(key)
	// Mark as deleted in the cache. It'll be deleted from tree on commit
	if cacheObj, ok := cache.storage[cacheKey]; ok {
		cache.storage[string(key)] = dataCacheObj{cacheObj.value, true, true}
	}
}

// Get a value for a given key.  Tries the cache first and then the state store
func (cache *KVCache) Get(key []byte) []byte {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()

	cacheKey := string(key)

	// check the cache
	data := cache.storage[cacheKey]
	if !data.delete && data.value != nil {
		return data.value
	}

	if data.value != nil && data.delete {
		return nil
	}

	// Not in the cache, go to cold storage
	value := cache.statedb.Get(key)
	if value != nil {
		cache.storage[cacheKey] = dataCacheObj{value, false, false}
		return value
	}

	return nil
}

// Iterate and produce results that are processed via the callback func
func (cache *KVCache) Iterate(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool {
	return cache.statedb.Iterate(start, end, ascending, fn)
}

// ApplyToState is called during abci.commit().   It sorts all keys in the cache
// for determinism, then writes the set to the tree.
func (cache *KVCache) ApplyToState() {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()

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

		if data.delete {
			cache.statedb.Delete([]byte(k))
			continue
		}
		cache.statedb.Set([]byte(k), data.value)
	}
}
