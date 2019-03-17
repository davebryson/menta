// Package app is the core menta application.  It contains logic to add callbacks, implements
// the ABCI interface, and provides and embedded node.
package app

import (
	"github.com/davebryson/menta/store"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
)

var _ abci.Application = (*MentaApp)(nil)

// MentaApp contains all the basics needed to build a tendermint application
type MentaApp struct {
	name         string
	state        *store.StateStore
	deliverCache *store.KVCache
	checkCache   *store.KVCache
	Config       *cfg.Config

	// tx decoder
	decodeTxFn sdk.TxDecoder
	// initChain
	initChainHandler sdk.InitialStartHandler
	// CheckTx
	checkTxHandler sdk.TxHandler
	// BeginBlock
	beginBlockHandler sdk.BeginBlockHandler
	// EndBlock
	endBlockHandler sdk.EndBlockHandler
	// DeliverTx router
	router map[string]sdk.TxHandler
	// Query router
	queryRouter map[string]store.QueryHandler
}

// NewApp returns a new instance of MentaApp where appname is the name of
// your application, and homedir is the path where menta/tendermint will
// store all the data and configuration information
// Should
func NewApp(appname, homedir string, txDecoder sdk.TxDecoder) *MentaApp {
	config, err := LoadConfig(homedir)
	if err != nil {
		panic(err)
	}

	state := store.NewStateStore(config.DBDir())
	app := &MentaApp{
		name:         appname,
		state:        state,
		Config:       config,
		deliverCache: state.RefreshCache().(*store.KVCache),
		checkCache:   state.RefreshCache().(*store.KVCache),
		router:       make(map[string]sdk.TxHandler, 0),
		queryRouter:  make(map[string]store.QueryHandler, 0),
	}
	app.decodeTxFn = txDecoder
	return app
}

// NewMockApp creates a menta app that can be used for local testing
// without a full blown node and an in memory state tree
func NewMockApp(txDecoder sdk.TxDecoder) *MentaApp {
	// Returns a inmemory app without tendermint for testing
	state := store.NewStateStore("")
	app := &MentaApp{
		name:         "mockapp",
		state:        state,
		deliverCache: state.RefreshCache().(*store.KVCache),
		checkCache:   state.RefreshCache().(*store.KVCache),
		router:       make(map[string]sdk.TxHandler, 0),
		queryRouter:  make(map[string]store.QueryHandler, 0),
	}
	app.decodeTxFn = txDecoder
	return app
}

// OnInitialStart : Add this handler if you'd like to do something
// the very *first* time Menta starts the new app.  Usually this is used
// to load initial state...accounts, etc....
func (app *MentaApp) OnInitialStart(fn sdk.InitialStartHandler) {
	app.initChainHandler = fn
}

// OnVerifyTx : Add this handler to validate your transactions.  This is NOT
// required as you can also verify a tx in your handler if you want.
// Usually this is a good place to put signature verification.  Txs that pass
// here by returing a result.ok will be added to the mempool for consideration
// in a block and them processed by your handlers registered with Route()
func (app *MentaApp) OnVerifyTx(fn sdk.TxHandler) {
	app.checkTxHandler = fn
}

// OnBeginBlock : Add this handler if you want to do something before
// a block of transactions are processed by your handler
func (app *MentaApp) OnBeginBlock(fn sdk.BeginBlockHandler) {
	app.beginBlockHandler = fn
}

// Route : Tells menta which handlers you want to run.
// This is the core 'business logic' for your application.
// 'routeName' should correspond with the value returned for a given msg.Route().
// 'fn' is your handler.
func (app *MentaApp) Route(routeName string, fn sdk.TxHandler) {
	app.router[routeName] = fn
}

func (app *MentaApp) RouteQuery(path string, fn store.QueryHandler) {
	app.queryRouter[path] = fn
}

// OnEndBlock : Add this handler to do something after a block of
// transactions are processed.  This is commonly used to update network
// validators based on logic implemented via an OnTx handler.
func (app *MentaApp) OnEndBlock(fn sdk.EndBlockHandler) {
	app.endBlockHandler = fn
}

// ---------------------------------------------------------------
//
// ABCI Implementations
//
// ---------------------------------------------------------------

// InitChain (ABCI callback) is ran once on the very first run of the application chain.
func (app *MentaApp) InitChain(req abci.RequestInitChain) (resp abci.ResponseInitChain) {
	// TODO: add Validators to state
	if app.initChainHandler != nil {
		ctx := sdk.NewContext(app.deliverCache, nil)
		resp = app.initChainHandler(ctx, req)
	}
	return
}

// Info (ABCI callback) checks the application state against Tendermint on startup.
// If the last block height known by the application is less than what Tendermint says,
// then your application node will sync by replaying all transactions up to the current
// tendermint block height.
func (app *MentaApp) Info(req abci.RequestInfo) abci.ResponseInfo {
	tmversion := req.GetVersion()
	return abci.ResponseInfo{
		Data:             app.name,
		Version:          tmversion,
		LastBlockHeight:  app.state.CommitInfo.Version,
		LastBlockAppHash: app.state.CommitInfo.Hash,
	}
}

// Query (ABCI callback) *committed* state in the Tree
func (app *MentaApp) Query(req abci.RequestQuery) abci.ResponseQuery {
	// Get path and key
	// look up handler by path
	// send it key and context
	if req.Data == nil || len(req.Data) == 0 {
		res := abci.ResponseQuery{}
		res.Code = sdk.BadQuery
		res.Log = "Error: query requires a key"
		return res
	}
	queryKey := req.Data
	queryPath := req.Path

	handler := app.queryRouter[queryPath]
	if handler == nil {
		res := abci.ResponseQuery{}
		res.Code = sdk.BadQuery
		res.Log = "no query handler found"
		return res
	}

	ctx := store.NewQueryContext(app.state)
	return handler(queryKey, ctx)
}

// CheckTx (ABCI callback) populates the mempool. Transactions are ran through
// OnVerifyTx(). If they pass, they will be considered for inclusion in a block
// and processed via DeliverTx() which calls your handlers
func (app *MentaApp) CheckTx(raw []byte) abci.ResponseCheckTx {
	if app.checkTxHandler == nil {
		// Nothing to do!
		return abci.ResponseCheckTx{}
	}

	tx, err := app.decodeTxFn(raw)
	if err != nil {
		e := sdk.ErrorBadTx()
		return abci.ResponseCheckTx{Code: e.Code, Log: e.Log}
	}

	ctx := sdk.NewContext(app.checkCache, tx)
	result := app.checkTxHandler(ctx)
	return abci.ResponseCheckTx{
		Code: result.Code,
		Log:  result.Log,
		Data: result.Data,
		Tags: result.Tags,
	}
}

// BeginBlock (ABCI callback) signals the start of processing a batch of transaction via DeliverTx
func (app *MentaApp) BeginBlock(req abci.RequestBeginBlock) (resp abci.ResponseBeginBlock) {
	if app.beginBlockHandler != nil {
		ctx := sdk.NewContext(app.deliverCache, nil)
		resp = app.beginBlockHandler(ctx, req)
	}
	return
}

// DeliverTx (ABCI callback) is the heart of processing transactions leading to a state transistion.
// This is where the your application logic (handlers) are executed.
func (app *MentaApp) DeliverTx(raw []byte) abci.ResponseDeliverTx {
	tx, err := app.decodeTxFn(raw)
	if err != nil {
		e := sdk.ErrorBadTx()
		return abci.ResponseDeliverTx{Code: e.Code, Log: e.Log}
	}

	handler := app.router[tx.GetMsg().Route()]
	if handler == nil {
		e := sdk.ErrorNoHandler()
		return abci.ResponseDeliverTx{Code: e.Code, Log: e.Log}
	}

	ctx := sdk.NewContext(app.deliverCache, tx)
	result := handler(ctx)

	return abci.ResponseDeliverTx{
		Code: result.Code,
		Log:  result.Log,
		Data: result.Data,
		Tags: result.Tags,
	}
}

// EndBlock (ABCI callback) signals the end of a block of txs.
// TODO: return changes to the validator set
func (app *MentaApp) EndBlock(req abci.RequestEndBlock) (resp abci.ResponseEndBlock) {
	if app.endBlockHandler != nil {
		ctx := sdk.NewContext(app.deliverCache, nil)
		resp = app.endBlockHandler(ctx, req)
	}
	return
}

// Commit (ABCI callback) : sends updates to the state tree and refresh caches
func (app *MentaApp) Commit() abci.ResponseCommit {
	app.deliverCache.ApplyToState()

	commitresults := app.state.Commit()

	app.deliverCache = app.state.RefreshCache().(*store.KVCache)
	app.checkCache = app.state.RefreshCache().(*store.KVCache)

	return abci.ResponseCommit{Data: commitresults.Hash}
}

// SetOption - (ABCI callback) not used
func (app *MentaApp) SetOption(req abci.RequestSetOption) abci.ResponseSetOption {
	return abci.ResponseSetOption{}
}
