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
	assert.Equal([]byte("dave"), dcache.Get([]byte("name")))
	assert.Nil(dcache.Get([]byte("not")))

	// abci.Commit()
	dcache.ApplyToState()
	info := st.Commit()
	assert.Equal(int64(1), info.Version)
	st.Close()

	// Check the store from the previous commit
	st = NewStateStore(".")
	dcache = NewCache(st)
	assert.Equal([]byte("dave"), dcache.Get([]byte("name")))
	assert.Equal(info.Version, st.CommitInfo.Version)
	assert.Equal(info.Hash, st.CommitInfo.Hash)

	dcache.Delete([]byte("name"))
	//dcache.ApplyToState()

	assert.Equal([]byte(nil), dcache.Get([]byte("name")))

	st.Close()
}

func TestStoreQueries(t *testing.T) {
	assert := assert.New(t)

	st := NewStateStore("") // in-memory
	dcache := NewCache(st)
	dcache.Set([]byte("name1"), []byte("dave"))
	dcache.Set([]byte("name2"), []byte("bob"))
	dcache.ApplyToState()
	st.Commit()

	rq1value := st.Get([]byte("name1"))
	assert.Equal([]byte("dave"), rq1value)

	//root := commit.Hash

	// Try proof
	//rqvalue, proof, err := st.Query([]byte("name1"), 0, true)
	//assert.Nil(err)
	//assert.NotNil(proof)
	//assert.NotNil(rqvalue)
	//proof.Ops[0].

	//rq2 := st.Query(abci.RequestQuery{Path: "/key", Data: []byte("name1"), Prove: true})
	//assert.NotNil(rq2)
	//assert.NotNil(rq2.Proof)
	// TODO: Verify the proof
	//fmt.Printf("Proof: %v\n", rq2.Proof)
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

	cache.Iterate([]byte("g1/"), []byte("g2/"), true, viewer)
	assert.Equal(4, len(viewed))
	assert.Equal("g1/s1", viewed[0])
	assert.Equal("g1/s4", viewed[3])

	allgs := []string{}
	cache.Iterate([]byte("g"), []byte("h"), true, func(key []byte, value []byte) bool {
		allgs = append(allgs, string(key))
		return false
	})
	assert.Equal(9, len(allgs))

}
