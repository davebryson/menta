package store

import (
	"sort"

	"github.com/davebryson/menta/types"
)

var _ types.Cache = (*KVCache)(nil)

type accountCacheObject struct {
	account *types.Account
	dirty   bool
}

type dataCacheObj struct {
	value []byte
	dirty bool
}

type KVCache struct {
	statedb  *StateStore
	storage  map[string]dataCacheObj
	accounts map[string]accountCacheObject
}

func NewCache(store *StateStore) *KVCache {
	return &KVCache{
		statedb:  store,
		storage:  make(map[string]dataCacheObj),
		accounts: make(map[string]accountCacheObject),
	}
}

func (cache *KVCache) SetAccount(account *types.Account) {
	key := cacheAccountKey(account.Address())
	cache.accounts[key] = accountCacheObject{account, true}
}

func (cache *KVCache) Set(key, val []byte) {
	cache.storage[string(key)] = dataCacheObj{val, true}
}

func (cache *KVCache) GetAccount(address []byte) (*types.Account, error) {
	key := cacheAccountKey(address)

	// Try the cache
	if acctObj, ok := cache.accounts[key]; ok {
		return acctObj.account, nil
	}

	// Not in the cache, go to the tree...
	acctBits := cache.statedb.Get([]byte(key))
	account, err := types.AccountFromBytes(acctBits)
	if err != nil {
		return nil, err
	}

	// Cache it
	cache.accounts[key] = accountCacheObject{account, false}

	return account, nil
}

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

func (cache *KVCache) ApplyToState() {
	// deterministic: sort storage and update in order

	// Sort account and update in order
	accountAddresses := make([]string, 0, len(cache.accounts))
	for k, _ := range cache.accounts {
		accountAddresses = append(accountAddresses, k)
	}

	sort.Strings(accountAddresses)

	for _, addy := range accountAddresses {
		a := cache.accounts[addy]
		if a.account == nil || !a.dirty {
			continue
		}
		account := a.account
		if bits, err := account.Bytes(); err == nil {
			cache.statedb.Set([]byte(addy), bits)
		}
	}

	storageKeys := make([]string, 0, len(cache.storage))
	for k, _ := range cache.storage {
		storageKeys = append(storageKeys, k)
	}

	sort.Strings(storageKeys)

	for _, k := range storageKeys {
		data := cache.storage[k]
		if data.value == nil || !data.dirty {
			continue
		}
		cache.statedb.Set([]byte(k), data.value)
	}
}
