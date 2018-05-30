package test

import (
	"os"
	"testing"

	"github.com/davebryson/menta/store"
	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
)

// Tests
func TestStoreBasics(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		os.RemoveAll("mstate.db")
	}()

	st := store.NewStateStore(".")
	dcache := store.NewCache(st)
	dcache.Set([]byte("name"), []byte("dave"))
	r := dcache.Get([]byte("name"))
	assert.Equal([]byte("dave"), r)
	r1 := dcache.Get([]byte("not"))
	assert.Nil(r1)

	// abci.Commit()
	dcache.ApplyToState()
	info := st.Commit()

	assert.Equal(int64(1), info.Version)
	st.Close()

	// Check the store
	st = store.NewStateStore(".")
	dcache = store.NewCache(st)
	r = dcache.Get([]byte("name"))

	assert.Equal([]byte("dave"), r)
	assert.Equal(info.Version, st.CommitInfo.Version)
	assert.Equal(info.Hash, st.CommitInfo.Hash)

	st.Close()
}

func TestStoreAndQuery(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		os.RemoveAll("mstate.db")
	}()

	k := sdk.CreateKey()
	bobAddress := k.Address
	bob, err := sdk.AccountFromPubKey(k.PrivateKey.PubKey().Bytes())
	assert.Nil(err)

	// Set up the store
	st := store.NewStateStore(".")
	cache := store.NewCache(st)
	cache.SetAccount(bob)
	cache.Set([]byte("t"), []byte("one"))
	// Commit
	cache.ApplyToState()
	c1 := st.Commit()

	// New Cache
	cache = store.NewCache(st)
	bob2, err := cache.GetAccount(bobAddress)
	assert.Nil(err)
	assert.Equal(bob.PubKey.Bytes(), bob2.PubKey.Bytes())
	result := cache.Get([]byte("t"))
	assert.Equal([]byte("one"), result)
	assert.Equal(st.CommitInfo.Version, c1.Version)

	// No changes
	cache.ApplyToState()
	c2 := st.Commit()
	assert.Equal(c2.Hash, c1.Hash)

	// New Cache
	cache = store.NewCache(st)
	// Add
	cache.Set([]byte("t"), []byte("two"))
	// Commit
	cache.ApplyToState()
	c3 := st.Commit()
	assert.NotEqual(c3.Hash, c2.Hash)
}
