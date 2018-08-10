package test

import (
	"fmt"
	"os"
	"sort"
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

func TestIter(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		os.RemoveAll("mstate.db")
	}()

	// From IAVL tests...
	type record struct {
		key   string
		value string
	}

	records := []record{
		{"abc", "123"},
		{"low", "high"},
		{"fan", "456"},
		{"foo", "a"},
		{"foobaz", "c"},
		{"good", "bye"},
		{"foobang", "d"},
		{"foobar", "b"},
		{"food", "e"},
		{"foml", "f"},
		{"g1/s1", "gs1"},
		{"g1/s2", "gs1"},
		{"g1/s3", "gs1"},
		{"g1/s4", "gs1"},
		{"g2/s1", "gs1"},
		{"g2/s2", "gs1"},
		{"g2/s3", "gs1"},
		{"g2/s4", "gs1"},
	}
	keys := make([]string, len(records))
	for i, r := range records {
		keys[i] = r.key
	}
	sort.Strings(keys)

	st := store.NewStateStore(".")
	cache := store.NewCache(st)
	for _, r := range records {
		cache.Set([]byte(r.key), []byte(r.value))
	}
	// Commit
	cache.ApplyToState()
	st.Commit()

	viewed := []string{}
	viewer := func(key []byte, value []byte) bool {
		viewed = append(viewed, string(key))
		return false
	}
	st.IterateKeyRange([]byte("g1/"), []byte("g2/"), true, viewer)
	assert.Equal(4, len(viewed))
	assert.Equal("g1/s1", viewed[0])
	assert.Equal("g1/s4", viewed[3])

	allgs := []string{}
	st.IterateKeyRange([]byte("g"), []byte("h"), true, func(key []byte, value []byte) bool {
		fmt.Printf("Key: %s\n", key)
		allgs = append(allgs, string(key))
		return false
	})
	assert.Equal(9, len(allgs))

}
