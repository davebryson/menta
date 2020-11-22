package storage

import "github.com/cosmos/iavl"

var _ TreeReader = (*Snapshot)(nil)

// Snapshot provides a read-only view of committed data
type Snapshot struct {
	// TODO: Use Immutable tree
	tree *iavl.MutableTree
}

// NewSnapshot is created via store.Snapshot()
func NewSnapshot(tree *iavl.MutableTree) Snapshot {
	return Snapshot{
		tree: tree,
	}
}

// GetWithProof returns proof of the key in the tree
func (snap Snapshot) GetWithProof(key []byte) ([]byte, *iavl.RangeProof, error) {
	return snap.tree.GetWithProof(key)
}

// Get a value
func (snap Snapshot) Get(key []byte) ([]byte, error) {
	_, bits := snap.tree.Get(key)
	if bits == nil {
		return nil, ErrValueNotFound
	}
	return bits, nil
}

// IterateKeyRange from start to end
func (snap Snapshot) IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool {
	return snap.tree.IterateRange(start, end, ascending, fn)
}
