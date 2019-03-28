package store

import (
	"fmt"

	"github.com/tendermint/go-amino"

	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"

	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	CacheSize   = 1000
	StateDbName = "mstate"
)

var commitKey = []byte("/menta/commitinfo")
var _ sdk.Store = (*StateStore)(nil)
var cdc = amino.NewCodec()

type StateStore struct {
	db         dbm.DB
	tree       *iavl.MutableTree
	CommitInfo *sdk.CommitInfo
	numHistory int64
}

func NewStateStore(dbdir string) *StateStore {
	db := loadDb(dbdir)
	ci := loadCommitData(db)
	tree := iavl.NewMutableTree(db, CacheSize)
	//fmt.Printf("Version here %v\n", ci)
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

// IterateKeyRange - iterator non-inclusive
func (st *StateStore) IterateKeyRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) bool {
	return st.tree.IterateRange(start, end, ascending, fn)
}

// Select a query version from the store.
func getQueryHeight(store *StateStore, queryVersion int64) int64 {
	tree := store.tree
	if queryVersion == 0 {
		latest := tree.Version()
		if tree.VersionExists(latest - 1) {
			// Use the last version vs the current version
			// if the user specifies 0 or makes no specfic request
			return latest - 1
		}
		// If no previous version exists, use the latest
		return latest
	}
	// Otherwise use what they requested
	return queryVersion
}

// Query returns a value and/or proof from the tree.
func (st *StateStore) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	// This code is all adapted from the Cosmos SDK
	if req.Data == nil || len(req.Data) == 0 {
		res.Code = sdk.BadQuery
		res.Log = "Error: query requires a key"
		return
	}

	queryKey := req.Data
	queryPath := req.Path
	queryVersion := getQueryHeight(st, req.Height)
	tree := st.tree

	switch queryPath {
	case "/key":
		res.Height = queryVersion
		res.Key = queryKey
		key := queryKey

		if req.Prove {
			value, proof, err := tree.GetVersionedWithProof(queryKey, queryVersion)
			if err != nil {
				res.Log = err.Error()
				break
			}
			if proof == nil {
				// Proof == nil implies that the store is empty.
				if value != nil {
					panic("unexpected value for an empty proof")
				}
			}
			if value != nil {
				res.Value = value
				res.Proof = &merkle.Proof{Ops: []merkle.ProofOp{iavl.NewIAVLValueOp(key, proof).ProofOp()}}
			} else {
				// value not found
				res.Value = nil
				res.Proof = &merkle.Proof{Ops: []merkle.ProofOp{iavl.NewIAVLAbsenceOp(key, proof).ProofOp()}}
			}
		} else {
			_, res.Value = tree.GetVersioned(queryKey, queryVersion)
		}
	default:
		res.Log = fmt.Sprintf("Unexpected Query path: %v", queryPath)
		res.Code = sdk.BadQuery
	}
	return
}

func (st *StateStore) Commit() *sdk.CommitInfo {
	hash, version, err := st.tree.SaveVersion()

	// from cosmos-sdk iavlstore - Release an old version of history
	if st.numHistory > 0 && (st.numHistory < st.tree.Version()) {
		toRelease := version - st.numHistory
		st.tree.DeleteVersion(toRelease)
	}

	latestCommit := &sdk.CommitInfo{Version: version, Hash: hash}
	bits, err := cdc.MarshalBinaryBare(latestCommit)
	if err != nil {
		panic(err)
	}

	// save commit to db
	st.db.Set(commitKey, bits)

	st.CommitInfo = latestCommit
	return latestCommit
}

func (st *StateStore) RefreshCache() sdk.Cache {
	return NewCache(st)
}

func (st *StateStore) Close() {
	st.db.Close()
}

func loadCommitData(db dbm.DB) *sdk.CommitInfo {
	commitBytes := db.Get(commitKey)
	if commitBytes != nil {
		ci := new(sdk.CommitInfo)
		err := cdc.UnmarshalBinaryBare(commitBytes, &ci)
		if err != nil {
			panic(err)
		}
		return ci
	}
	// Return a default
	return &sdk.CommitInfo{Version: 0, Hash: nil}
}

func loadDb(dbdir string) dbm.DB {
	if dbdir == "" {
		return dbm.NewMemDB()
	}
	return dbm.NewDB(StateDbName, dbm.GoLevelDBBackend, dbdir)
}
