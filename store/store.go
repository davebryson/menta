package store

import (
	sdk "github.com/davebryson/menta/types"
	proto "github.com/golang/protobuf/proto"

	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	cacheSize   = 1000
	stateDbName = "mstate"
	history     = int64(5)
)

var commitKey = []byte("/menta/commitinfo")

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
		numHistory: history,
	}
}

func (st *StateStore) Set(key, val []byte) {
	st.tree.Set(key, val)
}

func (st *StateStore) Delete(key []byte) {
	st.tree.Remove(key)
}

// IterateKeyRange - iterator non-inclusive
func (st *StateStore) Iterate(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool {
	return st.tree.IterateRange(start, end, ascending, fn)
}

func (st *StateStore) Get(key []byte) []byte {
	// This should call get...?
	_, bits := st.tree.Get(key)
	return bits
}

// NOT IN CURRENT USE - add back later.
/*func (st *StateStore) query(key []byte, requestedVersion int64, prove bool) ([]byte, *merkle.Proof, error) {
	// Default is the latest version of the tree
	version := st.tree.Version()

	if requestedVersion > 0 {
		if !st.tree.VersionExists(requestedVersion) {
			return nil, nil, fmt.Errorf("Store get value: requested version %d doesn't exist", requestedVersion)
		}
		version = requestedVersion
	}

	if !prove {
		_, result := st.tree.GetVersioned(key, version)
		return result, nil, nil
	}

	// They want proof...
	value, proof, err := st.tree.GetVersionedWithProof(key, version)
	if err != nil {
		return nil, nil, err
	}

	if proof == nil && value != nil {
		panic("Problem with state store!  Proof is nil for an existing value")
	}

	// Now return the appropriate proof depending on if there's a value of not
	if value == nil {
		// Return proof of non-existence
		proofNE := &merkle.Proof{Ops: []merkle.ProofOp{iavl.NewIAVLAbsenceOp(key, proof).ProofOp()}}
		return value, proofNE, nil
	}

	proofE := &merkle.Proof{Ops: []merkle.ProofOp{iavl.NewIAVLValueOp(key, proof).ProofOp()}}
	return value, proofE, nil
}*/

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

// RefreshCache on commit
func (st *StateStore) RefreshCache() sdk.Cache {
	return NewCache(st)
}

// Close the underlying database
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
	return dbm.NewDB(stateDbName, dbm.GoLevelDBBackend, dbdir)
}
