// Package app is the core menta application.  It provides a base Tendermint
// ABCI environment with pluggable Services.
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
	// CheckTx
	onValidationHandler sdk.ValidateTxHandler
	// DeliverTx router
	router map[string]sdk.Service
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
		router:       make(map[string]sdk.Service, 0),
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
		router:       make(map[string]sdk.Service, 0),
	}
}

// ValidateTxHandler : Add this handler to validate your transactions.  This is NOT
// required as you can also validate tx in your Service if you want.
// Usually this is a good place to put signature verification.
// There is only 1 of these per application.
func (app *MentaApp) ValidateTxHandler(fn sdk.ValidateTxHandler) {
	app.onValidationHandler = fn
}

// AddService : registers your service with Menta
func (app *MentaApp) AddService(service sdk.Service) {
	_, exists := app.router[service.Route()]
	if !exists {
		// First come, first serve
		app.router[service.Route()] = service
	}
}

// internal logic for check/deliverTx
func (app *MentaApp) runTx(rawtx []byte, isCheck bool) sdk.Result {
	tx, err := sdk.DecodeTx(rawtx)
	if err != nil {
		return sdk.ErrorBadTx()
	}

	_, ok := app.router[tx.Service]
	if !ok {
		return sdk.ResultError(1, "No service found!")
	}

	if isCheck {
		if app.onValidationHandler == nil {
			// Nothing to do!
			return sdk.Result{}
		}
		store := sdk.NewRWStore(tx.Service, app.checkCache)
		return app.onValidationHandler(tx, store)
	}

	store := sdk.NewRWStore(tx.Service, app.deliverCache)
	service := app.router[tx.Service]
	return service.Execute(tx, store)
}

// ---------------------------------------------------------------
//
// ABCI Callback Implementations
//
// ---------------------------------------------------------------

// InitChain is ran once, on the very first run of the application chain.
func (app *MentaApp) InitChain(req abci.RequestInitChain) (resp abci.ResponseInitChain) {
	for name, serv := range app.router {
		store := sdk.NewRWStore(name, app.deliverCache)
		serv.Init(store)
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
// This calls the handler where the path is the Service name return from
// Service.Route() and the key is the application specific key in storage
func (app *MentaApp) Query(query abci.RequestQuery) abci.ResponseQuery {
	res := abci.ResponseQuery{}
	if query.Data == nil || len(query.Data) == 0 {
		res.Code = sdk.BadQuery
		res.Log = "Error: query requires a key"
		return res
	}
	queryKey := query.Data
	serviceName := query.Path

	// Note: query
	service, ok := app.router[serviceName]
	if !ok {
		res.Code = sdk.BadQuery
		res.Log = "no query handler found"
		return res
	}

	store := sdk.NewQueryStore(service.Route(), app.deliverCache)
	result := service.Query(queryKey, store)

	res.Code = result.Code
	res.Value = result.Data
	res.Log = result.Log
	return res
}

// CheckTx populates the mempool. Transactions are ran through the OnValidationHandler.
// If the pass, they will be considered for inclusion in a block and processed via
// DeliverTx
func (app *MentaApp) CheckTx(checkTx abci.RequestCheckTx) abci.ResponseCheckTx {
	result := app.runTx(checkTx.Tx, true)
	return abci.ResponseCheckTx{
		Code: result.Code,
		Log:  result.Log,
		Data: result.Data,
	}
}

// BeginBlock signals the start of processing a batch of transaction via DeliverTx
func (app *MentaApp) BeginBlock(req abci.RequestBeginBlock) (resp abci.ResponseBeginBlock) {
	// Maybe add later
	return
}

// DeliverTx is the heart of processing transactions leading to a state transistion.
// This is where the your application logic lives via handlers
func (app *MentaApp) DeliverTx(dtx abci.RequestDeliverTx) abci.ResponseDeliverTx {
	result := app.runTx(dtx.Tx, false)
	return abci.ResponseDeliverTx{
		Code: result.Code,
		Log:  result.Log,
		Data: result.Data,
	}
}

// EndBlock signals the end of a block of txs.
// TODO: return changes to the validator set
func (app *MentaApp) EndBlock(req abci.RequestEndBlock) (resp abci.ResponseEndBlock) {
	// Add later
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
