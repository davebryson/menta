package types

// Callback handlers (functions) for Menta
import (
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitChainHandler callback to load genesis state
type InitChainHandler func(store RWStore, req abci.RequestInitChain) abci.ResponseInitChain

// TxHandler used for both check/deliver transactions
type TxHandler func(store RWStore, tx *Tx) Result

// QueryHandler used to query committed state from rpc calls
type QueryHandler func(store StoreReader, key []byte) ([]byte, error)

// BeginBlockHandler for any thing related to begin block events
type BeginBlockHandler func(store RWStore, req abci.RequestBeginBlock) abci.ResponseBeginBlock

// EndBlockHandler - do stuff at end block
type EndBlockHandler func(store RWStore, req abci.RequestEndBlock) abci.ResponseEndBlock
