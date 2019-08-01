package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Used for testing iterator
type indexValue struct {
	i uint64
	v []byte
}

func getList() List {
	// in-memory store
	store := NewCache(NewStateStore(""))
	return NewList(store, []byte("bob"))
}

func TestListPushPop(t *testing.T) {
	assert := assert.New(t)
	lst := getList()
	assert.Equal(uint64(0), lst.Len())
	assert.True(lst.IsEmpty())

	lst.Push([]byte("one"))
	lst.Push([]byte("two"))
	lst.Push([]byte("three"))

	assert.Equal(uint64(3), lst.Len())
	assert.False(lst.IsEmpty())

	v := lst.Pop()
	assert.Equal([]byte("three"), v)
	assert.Equal(uint64(2), lst.Len())
	assert.False(lst.IsEmpty())

	lst.Pop()
	lst.Pop()

	assert.Equal(uint64(0), lst.Len())
	assert.True(lst.IsEmpty())
}

func TestListGetSetTruncateExtendClear(t *testing.T) {
	assert := assert.New(t)
	lst := getList()

	_, err := lst.Get(0)
	assert.NotNil(err)
	assert.EqualError(err, "Index out of bounds")

	_, err = lst.Get(5)
	assert.EqualError(err, "Index out of bounds")

	lst.Push([]byte("one"))
	lst.Push([]byte("two"))
	lst.Push([]byte("three"))

	v, err := lst.Get(1)
	assert.Nil(err)
	assert.Equal([]byte("two"), v)

	err = lst.Set(uint64(5), []byte("nope"))
	assert.EqualError(err, "Index out of bounds")

	err = lst.Set(uint64(2), []byte("threechanged"))
	assert.Nil(err)

	v, err = lst.Get(2)
	assert.Nil(err)
	assert.Equal([]byte("threechanged"), v)
	assert.Equal(uint64(3), lst.Len())

	err = lst.Truncate(4)
	assert.NotNil(err)

	assert.Equal(uint64(3), lst.Len())

	err = lst.Truncate(1)
	assert.Nil(err)
	assert.Equal(uint64(1), lst.Len())

	newitems := [][]byte{
		[]byte("2"),
		[]byte("3"),
		[]byte("4"),
		[]byte("5"),
	}
	lst.Extend(newitems)
	assert.Equal(uint64(5), lst.Len())

	v, err = lst.Get(2)
	assert.Nil(err)
	assert.Equal([]byte("3"), v)

	err = lst.Clear()
	assert.Nil(err)
	assert.Equal(uint64(0), lst.Len())

	err = lst.Clear()
	assert.EqualError(err, "Index out of bounds")
}

func TestListIter(t *testing.T) {
	assert := assert.New(t)
	st := NewStateStore("")
	cache := NewCache(st)

	expected := []struct {
		i uint64
		v []byte
	}{
		{uint64(0), []byte("one")},
		{uint64(1), []byte("two")},
		{uint64(2), []byte("three")},
		{uint64(3), []byte("four")},
		{uint64(4), []byte("five")},
	}

	lst := NewList(cache, []byte("bob"))
	for _, item := range expected {
		lst.Push(item.v)
	}
	cache.ApplyToState() // MUST COMMIT. Iter reads from immutable tree!

	results := make([]indexValue, 0)
	lst.Iterate(func(i uint64, v []byte) (stop bool) {
		results = append(results, indexValue{i, v})
		return false
	})

	assert.Equal(len(expected), len(results))

	for ii, ex := range expected {
		assert.Equal(ex.i, results[ii].i)
		assert.Equal(ex.v, results[ii].v)
	}
}
