package types

import (
	abci "github.com/tendermint/abci/types"
)

type TxHandler func(ctx Context) Result

type GenesisHandler func(ctx Context, req abci.RequestInitChain)

type BeginBlockHandler func(ctx Context, req abci.RequestBeginBlock)

type EndBlockHandler func(ctx Context, req abci.RequestEndBlock) abci.ResponseEndBlock
