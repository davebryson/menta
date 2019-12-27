package storage

import "github.com/tendermint/iavl"

var _ TreeReader = (*Snapshot)(nil)

type Snapshot struct {
	tree *iavl.MutableTree
}

func NewSnaphot(tree *iavl.MutableTree) Snapshot {
	return Snapshot{
		tree: tree,
	}
}

func (snap Snapshot) Get(key []byte) ([]byte, error) {
	_, bits := snap.tree.Get(key)
	if bits == nil {
		return nil, ErrValueNotFound
	}
	return bits, nil
}

func (snap Snapshot) IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool {
	return snap.tree.IterateRange(start, end, ascending, fn)
}
