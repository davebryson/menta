package store

import abci "github.com/tendermint/tendermint/abci/types"

// QueryHandler for query routes
type QueryHandler func(key []byte, ctx QueryContext) abci.ResponseQuery

// QueryContext provides a read-only query context
type QueryContext struct {
	state *StateStore
}

// NewQueryContext creates an instance
func NewQueryContext(state *StateStore) QueryContext {
	return QueryContext{
		state: state,
	}
}

// Get a value for a given key
func (q QueryContext) Get(key []byte) []byte { return q.state.Get(key) }
