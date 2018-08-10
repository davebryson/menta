package store

import (
	"fmt"

	sdk "github.com/davebryson/menta/types"
	proto "github.com/golang/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tendermint/libs/db"
)

// Move to protobuf!

const (
	CacheSize   = 1000
	StateDbName = "mstate"
)

var commitKey = []byte("/menta/commitinfo")

//var cdc = amino.NewCodec()
var _ sdk.Store = (*StateStore)(nil)

type StateStore struct {
	db         dbm.DB
	tree       *iavl.VersionedTree
	CommitInfo sdk.CommitInfo
	numHistory int64
}

func NewStateStore(dbdir string) *StateStore {
	db := loadDb(dbdir)
	ci := loadCommitData(db)
	tree := iavl.NewVersionedTree(db, CacheSize)
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

func (st *StateStore) Get(key []byte) []byte {
	_, bits := st.tree.Get(key)
	return bits
}

// Basic iterator non-inclusive
func (st *StateStore) IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool {
	return st.tree.IterateRange(start, end, ascending, fn)
}

func (st *StateStore) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	// Based on the approach used in cosmos-sdk
	queryKey := req.Data
	queryPath := req.Path
	queryVersion := req.Height

	tree := st.tree
	if queryVersion == 0 {
		latest := tree.Version64()
		queryVersion = latest
		/*if tree.VersionExists(latest - 1) {
			queryVersion = latest - 1
		} else {
			queryVersion = latest
		}*/
	}

	res.Height = queryVersion
	res.Key = queryKey

	switch queryPath {
	case "/store", "/key":
		if req.Prove {
			value, proof, err := tree.GetVersionedWithProof(queryKey, queryVersion)
			if err != nil {
				res.Log = err.Error()
				break
			}
			res.Value = value
			res.Proof = []byte(proof.String()) // TODO: Is this right???
		} else {
			_, res.Value = tree.GetVersioned(queryKey, queryVersion)
		}
	default:
		res.Log = fmt.Sprintf("Unexpected Query path: %v", queryPath)
		res.Code = sdk.NotFound
	}
	return
}

func (st *StateStore) Commit() sdk.CommitInfo {
	hash, version, err := st.tree.SaveVersion()

	// from cosmos-sdk iavlstore - Release an old version of history
	if st.numHistory > 0 && (st.numHistory < st.tree.Version64()) {
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
