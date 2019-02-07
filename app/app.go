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

	// initChain
	onInitialStartHandler sdk.InitialStartHandler
	// CheckTx
	onValidationHandler sdk.TxHandler
	// BeginBlock
	onBeginBlockHandler sdk.BeginBlockHandler
	// EndBlock
	onEndBlockHandler sdk.EndBlockHandler
	// DeliverTx
	router *router
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
		router:       NewRouter(),
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
		router:       NewRouter(),
	}
}

// OnInitialStart : Add this handler if you'd like to do 'something'
// the very first time Menta starts the new app.  Usually this is used
// to load initial state...accounts, etc....
func (app *MentaApp) OnInitialStart(fn sdk.InitialStartHandler) {
	app.onInitialStartHandler = fn
}

// OnValidateTx : Add this handler to validate your transactions.  This is NOT
// required as you can also validate tx in your OnTx handlers if you want.
// Usually this is a good place to put signature verification.
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
	app.router.Add(routeName, fn)
}

// OnEndBlock : Add this handler to do something after a block of
// transactions are processed.  This is commonly used to update network
// validators based on logic implemented via an OnTx handler.
func (app *MentaApp) OnEndBlock(fn sdk.EndBlockHandler) {
	app.onEndBlockHandler = fn
}

// ---------------------------------------------------------------
//
// ABCI Implementations
//
// ---------------------------------------------------------------

// InitChain is ran once on the very first run of the application chain.
func (app *MentaApp) InitChain(req abci.RequestInitChain) (resp abci.ResponseInitChain) {
	// TODO: add Validators to state
	if app.onInitialStartHandler != nil {
		ctx := sdk.NewContext(app.deliverCache, nil)
		resp = app.onInitialStartHandler(ctx, req)
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
func (app *MentaApp) Query(query abci.RequestQuery) abci.ResponseQuery {
	// Delegate to Store
	return app.state.Query(query)
}

// CheckTx populates the mempool. Transactions are ran through the OnValidationHandler.
// If the pass, they will be considered for inclusion in a block and processed via
// DeliverTx
func (app *MentaApp) CheckTx(raw []byte) abci.ResponseCheckTx {
	if app.onValidationHandler == nil {
		// Nothing to do!
		return abci.ResponseCheckTx{}
	}

	tx, err := sdk.TransactionFromBytes(raw)
	if err != nil {
		e := sdk.ErrorBadTx()
		return abci.ResponseCheckTx{Code: e.Code, Log: e.Log}
	}

	ctx := sdk.NewContext(app.checkCache, tx)
	result := app.onValidationHandler(ctx)
	return abci.ResponseCheckTx{
		Code: result.Code,
		Log:  result.Log,
		Data: result.Data,
		Tags: result.Tags,
	}
}

// BeginBlock signals the start of processing a batch of transaction via DeliverTx
func (app *MentaApp) BeginBlock(req abci.RequestBeginBlock) (resp abci.ResponseBeginBlock) {
	if app.onBeginBlockHandler != nil {
		ctx := sdk.NewContext(app.deliverCache, nil)
		resp = app.onBeginBlockHandler(ctx, req)
	}
	return
}

// DeliverTx is the heart of processing transactions leading to a state transistion.
// This is where the your application logic lives via handlers
func (app *MentaApp) DeliverTx(raw []byte) abci.ResponseDeliverTx {
	tx, err := sdk.TransactionFromBytes(raw)
	if err != nil {
		e := sdk.ErrorBadTx()
		return abci.ResponseDeliverTx{Code: e.Code, Log: e.Log}
	}

	handler := app.router.GetHandler(tx.Call)
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

// EndBlock signals the end of a block of txs.
// TODO: return changes to the validator set
func (app *MentaApp) EndBlock(req abci.RequestEndBlock) (resp abci.ResponseEndBlock) {
	if app.onEndBlockHandler != nil {
		ctx := sdk.NewContext(app.deliverCache, nil)
		resp = app.onEndBlockHandler(ctx, req)
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
