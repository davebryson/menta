package types

// Callback handlers (functions) for Menta
import (
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitialStartHandler is called once the first time the application is ran.
// This is the place to load initial state in your application as it is
// passed the store to commit to tree state
type InitialStartHandler func(ctx Context, req abci.RequestInitChain) abci.ResponseInitChain

// BeginBlockHandler is called before DeliverTx
type BeginBlockHandler func(ctx Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock

// TxHandler are for check/delivery transactions. Will be called 1 or more times per block
type TxHandler func(ctx Context) Result

// EndBlockHandler is called after all DeliverTxs have been ran
type EndBlockHandler func(ctx Context, req abci.RequestEndBlock) abci.ResponseEndBlock
