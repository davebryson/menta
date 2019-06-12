package types

// Context provides access to the state store and
// the current decoded Transaction
type Context struct {
	Db Cache
	Tx *Tx
}

// NewContext creates a new context (hence the name)
func NewContext(db Cache, tx *Tx) Context {
	return Context{
		Db: db,
		Tx: tx,
	}
}

// QueryContext provides a read-only query context
type QueryContext struct {
	state KVStore
}

// NewQueryContext creates an instance
func NewQueryContext(state KVStore) QueryContext {
	return QueryContext{
		state: state,
	}
}

// Get a value for a given key
func (q QueryContext) Get(key []byte) []byte { return q.state.Get(key) }
