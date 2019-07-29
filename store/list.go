package store

import (
	"encoding/binary"
	"errors"
	"fmt"

	sdk "github.com/davebryson/menta/types"
)

var (
	// ErrOutOfBounds for the list
	ErrOutOfBounds = errors.New("Index out of bounds")
)

// List provides an array like structure over a key/value store using composite keys.
// The following structure is used in the underlying store:
//   /{key}/count => num     (where num is the total number of items in a list for a given key)
//   /{key}/{index} => value (where index is list index pointing to the value )
// Example:
//   Say our app requires the following: 'Each user has-many public keys'
//   We can use the List to capture:
//   'bob' is our key...maybe an account address. So,
//   /bob/count => 3 (this means bob has 3 public keys)
//   /bob/0 => publickey{}
//   /bob/1 => publickey{}
//   /bob/2 => publickey{}
// Everytime a new publickey{} is added for bob, the count is incremented.  We can use the
// count as an index into the list.
type List struct {
	key []byte
}

// NewList creates a new list for the given key.
// The key should be unique for the values its maintaining a list over
func NewList(key []byte) List {
	return List{
		key: key,
	}
}

// Len returns the current number of items in the list
func (list List) Len(store sdk.StoreReader) uint8 {
	value := store.Get(list.countKey())
	return decodeCountValue(value)
}

// IsEmpty - true if the list is empty
func (list List) IsEmpty(store sdk.StoreReader, key []byte) bool {
	len := list.Len(store)
	if len == 0 {
		return true
	}
	return false
}

// Checks if a given index is out of bounds
func (list List) isOutOfBounds(store sdk.StoreReader, index uint8) bool {
	if index+1 > list.Len(store) {
		return true
	}
	return false
}

// Get a value for a given key/index. The list uses zero-based index.
// Will return an error if the index is out of bounds
func (list List) Get(store sdk.StoreReader, index uint8) ([]byte, error) {
	if list.isOutOfBounds(store, index) {
		return nil, ErrOutOfBounds
	}
	return store.Get(list.indexKey(index)), nil
}

// Iterate over the list
func (list List) Iterate(store sdk.StoreReader, fn func(k []byte, v []byte) bool) bool {
	len := list.Len(store)
	startKey := list.indexKey(0)
	endKey := list.indexKey(len)
	return store.Iterate(startKey, endKey, true, fn)
}

// Append a value to the list
func (list List) Append(store sdk.RWStore, value []byte) {
	currentIndex := list.Len(store)
	store.Set(list.indexKey(currentIndex), value)
	store.Set(list.countKey(), encodeCountValue(currentIndex+1))
}

// Insert a value at a given index. This obviously will over-write the
// current value. It will return the previous value or an error if out of bounds
func (list List) Insert(store sdk.RWStore, index uint8, value []byte) ([]byte, error) {
	lastvalue, err := list.Get(store, index)
	if err != nil {
		return nil, err
	}
	store.Set(list.indexKey(index), value)
	return lastvalue, nil
}

// Pop (remove) an item at a given index.  This will delete the value,
// update all indices, and update the total count. Note, this can lead
// to 'write amplification' as 1 delete may translate to N+1 writes in
// order to update the list index.
func (list List) Pop(store sdk.RWStore, index uint8) error {
	if list.isOutOfBounds(store, index) {
		return ErrOutOfBounds
	}
	currentLen := list.Len(store)

	// Case 1: only 1 item in the list
	if currentLen == 1 {
		// Only delete the item at index 0
		store.Delete(list.indexKey(index))
		store.Set(list.countKey(), encodeCountValue(0))
		return nil
	}

	// Case 2: Delete the last item in the list
	if index == currentLen-1 {
		store.Delete(list.indexKey(index))
		store.Set(list.countKey(), encodeCountValue(currentLen-1))
		return nil
	}

	// Case 3: all others...
	newLength := currentLen - 1
	// Remove the item of interest
	store.Delete(list.indexKey(index))

	// Start from current index and rewrite the others
	for i := int(index); i < int(newLength); i++ {
		val, e := list.Get(store, uint8(i+1))
		if e != nil {
			return e
		}
		store.Set(list.indexKey(uint8(i)), val)
	}

	// Update the count
	store.Set(list.countKey(), encodeCountValue(newLength))
	return nil
}

// Key used to store the count for the list
func (list List) countKey() []byte {
	ck := fmt.Sprintf("/%s/count", list.key)
	return []byte(ck)
}

// Composite key for values
func (list List) indexKey(index uint8) []byte {
	ik := fmt.Sprintf("/%s/%04d", list.key, index)
	return []byte(ik)
}

// Serde for count value
func decodeCountValue(v []byte) uint8 {
	if len(v) == 0 {
		return 0
	}
	return uint8(binary.LittleEndian.Uint16(v))
}
func encodeCountValue(v uint8) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(v))
	return b
}
