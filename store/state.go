package store

import (
	"errors"

	sdk "github.com/davebryson/menta/types"
	proto "github.com/golang/protobuf/proto"

	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	cacheSize   = 10000
	StateDbName = "mstate"
)

var (
	commitKey     = []byte("/menta/commitinfo")
	ValueNotFound = errors.New("Store get: nil value for given key")
)

var _ sdk.Store = (*StateStore)(nil)

type StateStore struct {
	db         dbm.DB
	tree       *iavl.MutableTree
	CommitInfo sdk.CommitInfo
	numHistory int64
}

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

func (st *StateStore) Set(key, val []byte) {
	st.tree.Set(key, val)
}

func (st *StateStore) Delete(key []byte) {
	st.tree.Remove(key)
}

// GetCommitted returns only committed data, nothing cached
func (st *StateStore) GetCommitted(key []byte) ([]byte, error) {
	return st.Get(key)
}

func (st *StateStore) Get(key []byte) ([]byte, error) {
	_, bits := st.tree.Get(key)
	if bits == nil {
		return nil, ValueNotFound
	}
	return bits, nil
}

// IterateKeyRange - iterator non-inclusive
func (st *StateStore) IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool {
	return st.tree.IterateRange(start, end, ascending, fn)
}

func (st *StateStore) Commit() sdk.CommitInfo {
	hash, version, err := st.tree.SaveVersion()

	// from cosmos-sdk iavlstore - Release an old version of history
	if st.numHistory > 0 && (st.numHistory < st.tree.Version()) {
		toRelease := version - st.numHistory
		st.tree.DeleteVersion(toRelease)
	}

	// save commit to db
	com := sdk.CommitInfo{Version: version, Hash: hash}
	bits, err := proto.Marshal(&com)
	if err != nil {
		panic(err)
	}
	st.db.Set(commitKey, bits)

	st.CommitInfo = com
	return com
}

func (st *StateStore) RefreshCache() sdk.Cache {
	return NewCache(st)
}

func (st *StateStore) Close() {
	st.db.Close()
}

func loadCommitData(db dbm.DB) sdk.CommitInfo {
	commitBytes := db.Get(commitKey)
	var ci sdk.CommitInfo
	if commitBytes != nil {
		err := proto.Unmarshal(commitBytes, &ci)
		if err != nil {
			panic(err)
		}
	}
	return ci
}

func loadDb(dbdir string) dbm.DB {
	if dbdir == "" {
		return dbm.NewMemDB()
	}
	return dbm.NewDB(StateDbName, dbm.GoLevelDBBackend, dbdir)
}
