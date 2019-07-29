// Package app is the core menta application.  It contains logic to add callbacks, implements
// the ABCI interface, and provides and embedded node.
package app

import (
	"fmt"

	"github.com/davebryson/menta/store"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
)

const (
	badQueryCode        = 10
	handlerNotFoundCode = 11
	txDecodeErrorCode   = 12
)

var _ abci.Application = (*MentaApp)(nil)

// MentaApp contains all the basics needed to build a tendermint application
type MentaApp struct {
	name         string
	state        *store.StateStore
	deliverCache *store.KVCache
	checkCache   *store.KVCache
	Config       *cfg.Config

	// initChain
	onInitChainHandler sdk.InitChainHandler
	// CheckTx
	onValidationHandler sdk.TxHandler
	// BeginBlock
	onBeginBlockHandler sdk.BeginBlockHandler
	// EndBlock
	onEndBlockHandler sdk.EndBlockHandler
	// DeliverTx router
	router map[string]sdk.TxHandler
	// Query router
	queryRouter map[string]sdk.QueryHandler
}

// NewApp returns a new instance of MentaApp where appname is the name of
// your application, and homedir is the path where menta/tendermint will
// store all the data and configuration information
func NewApp(appname, homedir string) *MentaApp {
	config, err := LoadConfig(homedir)
	if err != nil {
		panic(err)
	}

	state := store.NewStateStore(config.DBDir())
	return &MentaApp{
		name:         appname,
		state:        state,
		Config:       config,
		deliverCache: state.RefreshCache().(*store.KVCache),
		checkCache:   state.RefreshCache().(*store.KVCache),
		router:       make(map[string]sdk.TxHandler, 0),
		queryRouter:  make(map[string]sdk.QueryHandler, 0),
	}
}

// NewMockApp creates a menta app that can be used for local testing
// without a full blown node and an in memory state tree
func NewMockApp() *MentaApp {
	// Returns a inmemory app without tendermint for testing
	state := store.NewStateStore("")
	return &MentaApp{
		name:         "mockapp",
		state:        state,
		deliverCache: state.RefreshCache().(*store.KVCache),
		checkCache:   state.RefreshCache().(*store.KVCache),
		router:       make(map[string]sdk.TxHandler, 0),
		queryRouter:  make(map[string]sdk.QueryHandler, 0),
	}
}

// OnInitChain : Add this handler if you'd like to do 'something'
// the very first time Menta starts the new app.  Usually this is used
// to load initial application state: accounts, etc....
func (app *MentaApp) OnInitChain(fn sdk.InitChainHandler) {
	app.onInitChainHandler = fn
}

// OnValidateTx : Add this handler to validate your transactions.  This is NOT
// required as you can also validate tx in your OnTx handlers if you want.
// Usually this is a good place to put signature verification.
// There is only 1 of these per application.
func (app *MentaApp) OnValidateTx(fn sdk.TxHandler) {
	app.onValidationHandler = fn
}

// OnBeginBlock : Add this handler if you want to do something before
// a block of transactions are processed.
func (app *MentaApp) OnBeginBlock(fn sdk.BeginBlockHandler) {
	app.onBeginBlockHandler = fn
}

// OnTx : This is the heart of the application.  TxHandlers are the
// application business logic.  The 'routeName' maps to the route field
// in the transaction, used to select which handler to select for which route.
// Transactions also have an 'action' field the can be used to further filter
// logic to sub functions under a single handler.
func (app *MentaApp) OnTx(routeName string, fn sdk.TxHandler) {
	app.router[routeName] = fn
}

// OnQuery : Add one or more handler to query application state.
func (app *MentaApp) OnQuery(routeName string, fn sdk.QueryHandler) {
	app.queryRouter[routeName] = fn
}

// OnEndBlock : Add this handler to do something after a block of
// transactions are processed.  This is commonly used to update network
// validators based on logic implemented via an OnTx handler.
func (app *MentaApp) OnEndBlock(fn sdk.EndBlockHandler) {
	app.onEndBlockHandler = fn
}

// ---------------------------------------------------------------
//
// ABCI Callback Implementations
//
// ---------------------------------------------------------------

// InitChain is ran once on the very first run of the application chain.
func (app *MentaApp) InitChain(req abci.RequestInitChain) (resp abci.ResponseInitChain) {
	if app.onInitChainHandler != nil {
		resp = app.onInitChainHandler(app.deliverCache, req)
	}
	return
}

// Info checks the application state on startup. If the last block height known by the
// application is less than what tendermint says, then the application node will sync
// by replaying all transactions up to the current tendermint block height.
func (app *MentaApp) Info(req abci.RequestInfo) abci.ResponseInfo {
	tmversion := req.GetVersion()
	return abci.ResponseInfo{
		Data:             app.name,
		Version:          tmversion,
		LastBlockHeight:  app.state.CommitInfo.Version,
		LastBlockAppHash: app.state.CommitInfo.Hash,
	}
}

// Query *committed* state in the Tree
// This calls the handler where the path is the registed query path (handler)
// and the key is the application specific key in storage
func (app *MentaApp) Query(query abci.RequestQuery) (resp abci.ResponseQuery) {
	if query.Data == nil || len(query.Data) == 0 {
		resp.Code = badQueryCode
		resp.Log = "App Query: query requires a key"
		return resp
	}
	key := query.Data
	route := query.Path

	handler := app.queryRouter[route]
	if handler == nil {
		resp.Code = badQueryCode
		resp.Log = fmt.Sprintf("App Query: No query handler found for route %s", route)
		return resp
	}

	result, err := handler(app.state, key)
	if err != nil {
		resp.Code = badQueryCode
		resp.Log = err.Error()
		return resp
	}

	resp.Code = 0
	resp.Value = result
	return resp
}

func (app *MentaApp) runTx(txbits []byte, isCheck bool) sdk.Result {
	// Try to decode it
	tx, err := sdk.DecodeTx(txbits)
	if err != nil {
		return sdk.ResultError(txDecodeErrorCode, err.Error())
	}

	// Check for a valid handler
	handler := app.router[tx.GetRoute()]
	if handler == nil {
		log := fmt.Sprintf("App: handler not found for route %s", tx.GetRoute())
		return sdk.ResultError(handlerNotFoundCode, log)
	}

	// If checkTx run that...
	if isCheck {
		if app.onValidationHandler == nil {
			// Nothing to do!
			return sdk.Result{}
		}
		return app.onValidationHandler(app.checkCache, tx)
	}

	// handle deliverTx
	handler = app.router[tx.GetRoute()]
	return handler(app.deliverCache, tx)
}

// CheckTx populates the mempool. Transactions are ran through the OnValidationHandler.
// If the pass, they will be considered for inclusion in a block and processed via
// DeliverTx
func (app *MentaApp) CheckTx(req abci.RequestCheckTx) (resp abci.ResponseCheckTx) {
	result := app.runTx(req.Tx, true)
	resp.Code = result.Code
	resp.Log = result.Log
	return resp
}

// DeliverTx is the heart of processing transactions leading to a state transistion.
// This is where the your application logic lives via handlers
func (app *MentaApp) DeliverTx(req abci.RequestDeliverTx) (resp abci.ResponseDeliverTx) {
	result := app.runTx(req.Tx, false)
	resp.Code = result.Code
	resp.Log = result.Log
	return resp
}

// BeginBlock signals the start of processing a batch of transaction via DeliverTx
func (app *MentaApp) BeginBlock(req abci.RequestBeginBlock) (resp abci.ResponseBeginBlock) {
	if app.onBeginBlockHandler != nil {
		resp = app.onBeginBlockHandler(app.deliverCache, req)
	}
	return
}

// EndBlock signals the end of a block of txs.
// TODO: return changes to the validator set
func (app *MentaApp) EndBlock(req abci.RequestEndBlock) (resp abci.ResponseEndBlock) {
	if app.onEndBlockHandler != nil {
		resp = app.onEndBlockHandler(app.deliverCache, req)
	}
	return
}

// Commit to state tree, refresh caches
func (app *MentaApp) Commit() abci.ResponseCommit {
	app.deliverCache.ApplyToState()

	commitresults := app.state.Commit()

	app.deliverCache = app.state.RefreshCache().(*store.KVCache)
	app.checkCache = app.state.RefreshCache().(*store.KVCache)

	return abci.ResponseCommit{Data: commitresults.Hash}
}

// SetOption - not used
func (app *MentaApp) SetOption(req abci.RequestSetOption) abci.ResponseSetOption {
	return abci.ResponseSetOption{}
}
