package storage

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

	st := NewStore(".")
	dcache := NewCache(st.Snapshot())
	// Setter/getter
	dcache.Put([]byte("name"), []byte("dave"))
	val, err := dcache.Get([]byte("name"))
	assert.Nil(err)
	assert.Equal([]byte("dave"), val)
	assert.Nil(dcache.Get([]byte("not")))

	// abci.Commit()
	info := st.Commit(dcache.ToBatch())
	assert.Equal(int64(1), info.Version)
	st.Close()

	// Check the store from the previous commit
	st = NewStore(".")
	snapshot := st.Snapshot()
	dcache = NewCache(snapshot)
	val, err = snapshot.Get([]byte("name"))
	assert.Nil(err)
	assert.Equal([]byte("dave"), val)
	assert.Equal(info.Version, st.CommitInfo.Version)
	assert.Equal(info.Hash, st.CommitInfo.Hash)

	for _, val := range [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e")} {
		dcache.Put(val, val)
	}

	val, err = dcache.Get([]byte("c"))
	assert.Nil(err)
	assert.Equal([]byte("c"), val)

	dcache.Remove([]byte("d"))
	_, err = dcache.Get([]byte("d"))
	assert.NotNil(err)

	st.Commit(dcache.ToBatch())
	_, err = st.Snapshot().Get([]byte("d"))
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

	st := NewStore(".")
	cache := NewCache(st.Snapshot())
	for _, r := range records {
		cache.Put([]byte(r.key), []byte(r.value))
	}
	// Commit
	st.Commit(cache.ToBatch())

	viewed := []string{}
	viewer := func(key []byte, value []byte) bool {
		viewed = append(viewed, string(key))
		return false
	}
	snapshot := st.Snapshot()

	snapshot.IterateKeyRange([]byte("g1/"), []byte("g2/"), true, viewer)
	assert.Equal(4, len(viewed))
	assert.Equal("g1/s1", viewed[0])
	assert.Equal("g1/s4", viewed[3])

	allgs := []string{}
	snapshot.IterateKeyRange([]byte("g"), []byte("h"), true, func(key []byte, value []byte) bool {
		allgs = append(allgs, string(key))
		return false
	})
	assert.Equal(9, len(allgs))

}
