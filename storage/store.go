package storage

import (
	"errors"
	fmt "fmt"
	"sort"

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

var _ TreeWriter = (*Store)(nil)

// Store provides access the the levelDb and Tree
type Store struct {
	db         dbm.DB
	tree       *iavl.MutableTree
	CommitInfo CommitData
	numHistory int64
}

// NewStore creates a new instance.  If 'dbdir' == "", it'll
// return an in-memory database
func NewStore(dbdir string) *Store {
	db := loadDb(dbdir)
	ci := loadCommitData(db)
	tree := iavl.NewMutableTree(db, cacheSize)
	tree.LoadVersion(ci.Version)

	return &Store{
		db:         db,
		tree:       tree,
		CommitInfo: ci,
		numHistory: 2, // Arbitrary for now...
	}
}

func (st *Store) Snapshot() TreeReader {
	// Note: Could use immutable tree here. But not sure how have that operation is
	return NewSnaphot(st.tree)
}

// Commit information about the current state to storage
func (st *Store) Commit(batch map[string]CacheOp) CommitData {
	fmt.Println("COMMIT")
	storageKeys := make([]string, 0, len(batch))
	for key := range batch {
		storageKeys = append(storageKeys, key)
	}
	// Sort keys for determinism (required by IAVL)
	sort.Strings(storageKeys)

	for _, key := range storageKeys {
		data := batch[key]
		// do delete and continue
		if data.delete {
			st.tree.Remove([]byte(key))
			continue
		}
		// Only insert dirty data. We don't re-insert unchanged, cached data
		if data.dirty {
			st.tree.Set([]byte(key), data.value)
		}
	}

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
