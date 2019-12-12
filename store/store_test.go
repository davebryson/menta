package store

import (
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests Store and Cache
func TestStoreBasics(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		os.RemoveAll("mstate.db")
	}()

	st := NewStateStore(".")
	dcache := NewCache(st)
	// Setter/getter
	dcache.Set([]byte("name"), []byte("dave"))
	val, err := dcache.Get([]byte("name"))
	assert.Nil(err)
	assert.Equal([]byte("dave"), val)
	assert.Nil(dcache.Get([]byte("not")))

	// abci.Commit()
	dcache.ApplyToState()
	info := st.Commit()
	assert.Equal(int64(1), info.Version)
	st.Close()

	// Check the store from the previous commit
	st = NewStateStore(".")
	dcache = NewCache(st)
	val, err = dcache.Get([]byte("name"))
	assert.Nil(err)
	assert.Equal([]byte("dave"), val)
	assert.Equal(info.Version, st.CommitInfo.Version)
	assert.Equal(info.Hash, st.CommitInfo.Hash)

	for _, val := range [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e")} {
		dcache.Set(val, val)
	}

	val, err = dcache.Get([]byte("c"))
	assert.Nil(err)
	assert.Equal([]byte("c"), val)

	dcache.Delete([]byte("d"))
	_, err = dcache.Get([]byte("d"))
	assert.NotNil(err)

	dcache.ApplyToState()
	_, err = dcache.Get([]byte("d"))
	assert.NotNil(err)

	st.Close()
}

func TestStoreIter(t *testing.T) {
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

	st := NewStateStore(".")
	cache := NewCache(st)
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

	cache.IterateKeyRange([]byte("g1/"), []byte("g2/"), true, viewer)
	assert.Equal(4, len(viewed))
	assert.Equal("g1/s1", viewed[0])
	assert.Equal("g1/s4", viewed[3])

	allgs := []string{}
	cache.IterateKeyRange([]byte("g"), []byte("h"), true, func(key []byte, value []byte) bool {
		allgs = append(allgs, string(key))
		return false
	})
	assert.Equal(9, len(allgs))

}
