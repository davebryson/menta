package storage

var _ Cache = (*KVCache)(nil)

type CacheOp struct {
	value  []byte
	dirty  bool
	delete bool
}

// KVCache provides a cached used for r/w access to storage
type KVCache struct {
	snapshot TreeReader
	storage  map[string]CacheOp
}

// NewCache return a fresh empty cache with ref to the State Store
func NewCache(snap TreeReader) *KVCache {
	return &KVCache{
		snapshot: snap,
		storage:  make(map[string]CacheOp),
	}
}

// Put a key in the cache
func (cache *KVCache) Put(key, val []byte) {
	cache.storage[string(key)] = CacheOp{val, true, false}
}

// Remove a key/value
func (cache *KVCache) Remove(key []byte) {
	if cache.Has(key) {
		cache.storage[string(key)] = CacheOp{nil, false, true}
	}
}

// Has - checks for a given key
func (cache *KVCache) Has(key []byte) bool {
	_, err := cache.Get(key)
	if err == nil {
		return true
	}
	return false
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
	value, err := cache.snapshot.Get(key)
	if err == nil {
		// cache it as not-dirty
		cache.storage[cacheKey] = CacheOp{value, false, false}
		return value, nil
	}

	return nil, ErrValueNotFound
}

// ToBatch returns the cached entries
func (cache *KVCache) ToBatch() map[string]CacheOp {
	return cache.storage
}
