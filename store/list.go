package store

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	// ErrOutOfBounds for the list
	ErrOutOfBounds = errors.New("Index out of bounds")
)

// List provides an array like structure over a key/value store using composite keys.
// The following format is used in the realize a list:
//   /{key}/count => num     (where num is the total number of items in a list for a given key)
//   /{key}/{index} => value (where index is list index pointing to the value )
// key:  Is a unique key used across the list
// Example:
//   Say we want to model that 'Each user has-many public keys'
//   We can use the List to capture this:
//   'bob' is our key...maybe an account address. So,
//   /bob/count => 3 (this means bob has 3 public keys)
//   /bob/0 => publickey{}
//   /bob/1 => publickey{}
//   /bob/2 => publickey{}
// Everytime a new publickey{} is added for bob, the count is incremented and we use the
// count as an index into the list.
type List struct {
	key   []byte
	store *KVCache
}

// NewList returns a new or existing list based on the key
func NewList(st *KVCache, key []byte) List {
	return List{
		key:   key,
		store: st,
	}
}

// Len returns the current length
func (l List) Len() uint64 {
	value := l.store.Get(l.countKey())
	return decode(value)
}

// internal
func (l List) isOutOfBounds(index uint64) bool {
	return index+1 > l.Len()
}

// IsEmpty -
func (l List) IsEmpty() bool {
	return l.Len() == 0
}

// Push a new value on to the list
func (l List) Push(value []byte) {
	cap := l.Len()
	l.store.Set(l.indexKey(cap), value)
	l.store.Set(l.countKey(), encode(cap+1))
}

// Pop removes and returns a value from the end of the list
func (l List) Pop() []byte {
	cap := l.Len()
	if cap == 0 {
		return nil
	}
	index := cap - 1
	key := l.indexKey(index)
	value := l.store.Get(key)
	l.store.Delete(key)
	l.store.Set(l.countKey(), encode(index))
	return value
}

// Set (overwrite) a value at a given index.  This will return an
// error if the index is out of bounds
func (l List) Set(index uint64, value []byte) error {
	if l.isOutOfBounds(index) {
		return ErrOutOfBounds
	}
	l.store.Set(l.indexKey(index), value)
	return nil
}

// Extend the list for the given values, increasing the len
// Ex:  If list A = [1,2] - A.extend([3,4,5]) => [1,2,3,4,5]
func (l List) Extend(values [][]byte) {
	for _, v := range values {
		l.Push(v)
	}
}

// Truncate - 'chop' off the end of the list at a given index
// Ex:  If list A =  [1,2,3,4,5]  A.Truncate(2) => [1,2]
func (l List) Truncate(index uint64) error {
	if l.isOutOfBounds(index) {
		return ErrOutOfBounds
	}
	totalLength := l.Len()
	count := totalLength
	for i := index; i < totalLength; i++ {
		l.store.Delete(l.indexKey(i))
		count--
		l.store.Set(l.countKey(), encode(count))
	}
	return nil
}

// Clear the list removing all values
func (l List) Clear() error {
	return l.Truncate(0)
}

// Get a value at a given index. Will return an error
// if the given index is out of bounds.
func (l List) Get(index uint64) ([]byte, error) {
	if l.isOutOfBounds(index) {
		return nil, ErrOutOfBounds
	}
	return l.store.Get(l.indexKey(index)), nil
}

// Iterate over entries in the *committed* list. The callback function will be passed
// each visited index and value. NOTE: The iterator will not return un-committed entries
// in the current cache
func (l List) Iterate(fn func(index uint64, value []byte) bool) bool {
	end := l.Len()
	index := -1
	return l.store.Iterate(l.indexKey(0), l.indexKey(end), true, func(_k, v []byte) bool {
		index++
		return fn(uint64(index), v)
	})
}

// internal - count key format
func (l List) countKey() []byte {
	return []byte(fmt.Sprintf("/%x/0x01", l.key))
}

// internal - Composite key format for values
func (l List) indexKey(index uint64) []byte {
	return []byte(fmt.Sprintf("/%x/%020d", l.key, index))
}

// internal - encode the count value as a []byte
func encode(count uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, count)
	return b
}

// internal - decode a []byte count to a uint64
func decode(raw []byte) uint64 {
	if raw == nil || len(raw) == 0 {
		return 0
	}
	return binary.LittleEndian.Uint64(raw)
}
