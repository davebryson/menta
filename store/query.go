package store

import abci "github.com/tendermint/tendermint/abci/types"

// Read-only query context.  Doesn't pay attention to tree versions
type QueryContext struct {
	state *StateStore
}

func NewQueryContext(state *StateStore) QueryContext {
	return QueryContext{
		state: state,
	}
}
func (q QueryContext) Get(key []byte) []byte { return q.state.Get(key) }

// Handler for query routes
type QueryHandler func(key []byte, ctx QueryContext) abci.ResponseQuery
