package store

import (
	"os"
	"sort"
	"testing"

	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Tests Store and Cache
func TestStoreBasics(t *testing.T) {
	sdk.RegisterStandardTypes(cdc)
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

	rq1 := st.Query(abci.RequestQuery{Path: "/key", Data: []byte("name1")})
	assert.NotNil(rq1)
	assert.Equal(uint32(0), rq1.Code)
	assert.Equal([]byte("dave"), rq1.Value)

	//root := commit.Hash

	// Try proof
	rq2 := st.Query(abci.RequestQuery{Path: "/key", Data: []byte("name1"), Prove: true})
	assert.NotNil(rq2)
	assert.NotNil(rq2.Proof)
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
