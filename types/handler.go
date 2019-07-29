package types

// Callback handlers (functions) for Menta
import (
	abci "github.com/tendermint/tendermint/abci/types"
)

type InitChainHandler func(store RWStore, req abci.RequestInitChain) abci.ResponseInitChain
type TxHandler func(store RWStore, tx *Tx) Result
type QueryHandler func(store StoreReader, key []byte) ([]byte, error)

type BeginBlockHandler func(store RWStore, req abci.RequestBeginBlock) abci.ResponseBeginBlock
type EndBlockHandler func(store RWStore, req abci.RequestEndBlock) abci.ResponseEndBlock
