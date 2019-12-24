package store

import (
	"errors"

	proto "github.com/golang/protobuf/proto"

	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	cacheSize = 10000
	// StateDbName is the filename of the kvstore
	StateDbName = "mstate"
)

var (
	commitKey = []byte("/menta/commitinfo")
	// ErrValueNotFound returned when the value for a key is nil
	ErrValueNotFound = errors.New("Store get: nil value for given key")
)

var _ Store = (*StateStore)(nil)

// StateStore provides access the the levelDb and Tree
type StateStore struct {
	db         dbm.DB
	tree       *iavl.MutableTree
	CommitInfo CommitData
	numHistory int64
}

// NewStateStore creates a new instance.  If 'dbdir' == "", it'll
// return an in-memory database
func NewStateStore(dbdir string) *StateStore {
	db := loadDb(dbdir)
	ci := loadCommitData(db)
	tree := iavl.NewMutableTree(db, cacheSize)
	tree.LoadVersion(ci.Version)

	return &StateStore{
		db:         db,
		tree:       tree,
		CommitInfo: ci,
		numHistory: 10, // Arbitrary for now...
	}
}

// Set a k/v in the tree
func (st *StateStore) Set(key, val []byte) {
	st.tree.Set(key, val)
}

// Delete a k/v from the tree
func (st *StateStore) Delete(key []byte) {
	st.tree.Remove(key)
}

// GetCommitted returns only committed data, nothing cached
// ** Implemented here to satisfy the KVStore interface.  Need
// to improve this
func (st *StateStore) GetCommitted(key []byte) ([]byte, error) {
	return st.Get(key)
}

// Get returns committed data from the tree
func (st *StateStore) Get(key []byte) ([]byte, error) {
	_, bits := st.tree.Get(key)
	if bits == nil {
		return nil, ErrValueNotFound
	}
	return bits, nil
}

func (st *StateStore) Has(key []byte) bool {
	_, err := st.Get(key)
	return err != nil
}

// IterateKeyRange - iterator non-inclusive
func (st *StateStore) IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool {
	return st.tree.IterateRange(start, end, ascending, fn)
}

// Commit information about the current state to the db
func (st *StateStore) Commit() CommitData {
	hash, version, err := st.tree.SaveVersion()

	// from cosmos-sdk iavlstore - Release an old version of history
	if st.numHistory > 0 && (st.numHistory < st.tree.Version()) {
		toRelease := version - st.numHistory
		st.tree.DeleteVersion(toRelease)
	}

	// save commit to db
	com := CommitData{Version: version, Hash: hash}
	bits, err := proto.Marshal(&com)
	if err != nil {
		panic(err)
	}
	st.db.Set(commitKey, bits)

	st.CommitInfo = com
	return com
}

// RefreshCache clears existing k/v from the cache
func (st *StateStore) RefreshCache() Cache {
	return NewCache(st)
}

// Close the underlying db
func (st *StateStore) Close() {
	st.db.Close()
}

// LoadCommitData from the db
func loadCommitData(db dbm.DB) CommitData {
	commitBytes := db.Get(commitKey)
	var ci CommitData
	if commitBytes != nil {
		err := proto.Unmarshal(commitBytes, &ci)
		if err != nil {
			panic(err)
		}
	}
	return ci
}

// load the db
func loadDb(dbdir string) dbm.DB {
	if dbdir == "" {
		return dbm.NewMemDB()
	}
	return dbm.NewDB(StateDbName, dbm.GoLevelDBBackend, dbdir)
}
