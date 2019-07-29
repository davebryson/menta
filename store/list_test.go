package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	assert := assert.New(t)
	st := NewStateStore("") // in-memory
	cache := NewCache(st)
	list := NewList([]byte("bob"))

	assert.Equal(uint8(0), list.Len(cache))

	list.Append(cache, []byte("one"))
	assert.Equal(uint8(1), list.Len(cache))
	v, e := list.Get(cache, uint8(0))
	assert.Nil(e)
	assert.Equal([]byte("one"), v)
	cache.ApplyToState()

	// Check again when committed
	assert.Equal(uint8(1), list.Len(cache))
	v, e = list.Get(cache, uint8(0))
	assert.Nil(e)
	assert.Equal([]byte("one"), v)

	list.Append(cache, []byte("two"))
	list.Append(cache, []byte("three"))
	list.Append(cache, []byte("four"))

	assert.Equal(uint8(4), list.Len(cache))
	v, e = list.Get(cache, uint8(2))
	assert.Equal([]byte("three"), v)
	cache.ApplyToState()

	// OutOfBounds
	_, e1 := list.Get(cache, 100)
	assert.NotNil(e1)

	results := make([][]byte, 0)
	list.Iterate(cache, func(k []byte, v []byte) bool {
		results = append(results, v)
		return false
	})

	assert.Equal(4, len(results))
	assert.Equal([]byte("one"), results[0])
}

func TestDeleteFromHead(t *testing.T) {
	assert := assert.New(t)
	st := NewStateStore("") // in-memory
	cache := NewCache(st)
	list := NewList([]byte("bob"))

	list.Append(cache, []byte("one"))
	list.Append(cache, []byte("two"))
	list.Append(cache, []byte("three"))

	list.Pop(cache, 0) // Delete 'one'
	v, e := list.Get(cache, 0)
	assert.Nil(e)
	assert.Equal([]byte("two"), v)

	v, e = list.Get(cache, 1)
	assert.Nil(e)
	assert.Equal([]byte("three"), v)
	assert.Equal(uint8(2), list.Len(cache))
}

func TestDeleteFromMiddle(t *testing.T) {
	assert := assert.New(t)
	st := NewStateStore("") // in-memory
	cache := NewCache(st)
	list := NewList([]byte("bob"))

	list.Append(cache, []byte("one"))
	list.Append(cache, []byte("two"))
	list.Append(cache, []byte("three"))
	list.Append(cache, []byte("four"))
	list.Append(cache, []byte("five"))

	list.Pop(cache, 0) // Delete 'three'
	v, e := list.Get(cache, 2)
	assert.Nil(e)
	assert.Equal([]byte("four"), v)

	v, e = list.Get(cache, 3)
	assert.Nil(e)
	assert.Equal([]byte("five"), v)

	list.Pop(cache, 1) // Delete 'two'
	v, e = list.Get(cache, 1)
	assert.Nil(e)
	assert.Equal([]byte("four"), v)

	assert.Equal(uint8(3), list.Len(cache))
}

func TestDeleteTail(t *testing.T) {
	assert := assert.New(t)
	st := NewStateStore("") // in-memory
	cache := NewCache(st)
	list := NewList([]byte("bob"))

	list.Append(cache, []byte("one"))
	list.Append(cache, []byte("two"))
	list.Append(cache, []byte("three"))
	list.Append(cache, []byte("four"))

	list.Pop(cache, 3) // Delete 'four'
	v, e := list.Get(cache, 0)
	assert.Nil(e)
	assert.Equal([]byte("one"), v)

	v, e = list.Get(cache, 1)
	assert.Nil(e)
	assert.Equal([]byte("two"), v)

	v, e = list.Get(cache, 2)
	assert.Nil(e)
	assert.Equal([]byte("three"), v)

	assert.Equal(uint8(3), list.Len(cache))
}
