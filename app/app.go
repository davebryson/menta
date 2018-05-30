package app

import (
	"github.com/davebryson/menta/store"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
)

var _ abci.Application = (*MentaApp)(nil)

type MentaApp struct {
	name         string
	state        *store.StateStore
	deliverCache *store.KVCache
	checkCache   *store.KVCache
	config       *cfg.Config

	// initChain
	onGenesisHandler sdk.GenesisHandler
	// CheckTx
	onValidationHandler sdk.TxHandler
	// BeginBlock
	onBeginBlockHandler sdk.BeginBlockHandler
	// EndBlock
	onEndBlockHandler sdk.EndBlockHandler
	// DeliverTx
	handlers map[string]sdk.TxHandler
}

func NewApp(appname, homedir string) *MentaApp {
	config, err := LoadConfig(homedir)
	if err != nil {
		panic(err)
	}

	state := store.NewStateStore(config.DBDir())
	return &MentaApp{
		name:         appname,
		state:        state,
		config:       config,
		deliverCache: state.RefreshCache().(*store.KVCache),
		checkCache:   state.RefreshCache().(*store.KVCache),
		handlers:     make(map[string]sdk.TxHandler),
	}
}

func NewTestApp() *MentaApp {
	// Returns a inmemory app without tendermint for testing
	state := store.NewStateStore("")
	return &MentaApp{
		name:         "testapp",
		state:        state,
		deliverCache: state.RefreshCache().(*store.KVCache),
		checkCache:   state.RefreshCache().(*store.KVCache),
		handlers:     make(map[string]sdk.TxHandler),
	}
}

func (app *MentaApp) OnGenesis(fn sdk.GenesisHandler) {
	app.onGenesisHandler = fn
}

func (app *MentaApp) OnValidateTx(fn sdk.TxHandler) {
	app.onValidationHandler = fn
}

func (app *MentaApp) OnBeginBlock(fn sdk.BeginBlockHandler) {
	app.onBeginBlockHandler = fn
}

func (app *MentaApp) OnTx(id string, fn sdk.TxHandler) {
	app.handlers[id] = fn
}

func (app *MentaApp) OnEndBlock(fn sdk.EndBlockHandler) {
	app.onEndBlockHandler = fn
}

// -----  ABCI  ------

// Check state on startup.  returned values determine sync
func (app *MentaApp) Info(req abci.RequestInfo) abci.ResponseInfo {
	return abci.ResponseInfo{
		Data:             app.name,
		LastBlockHeight:  app.state.CommitInfo.Version,
		LastBlockAppHash: app.state.CommitInfo.Hash,
	}
}

// Set application option
func (app *MentaApp) SetOption(req abci.RequestSetOption) abci.ResponseSetOption {
	return abci.ResponseSetOption{}
}

// Query State
func (app *MentaApp) Query(query abci.RequestQuery) abci.ResponseQuery {
	return app.state.Query(query)
}

// Mempool Connection - validate a transaction for inclusion in mempoo
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

// Consensus Connections

// Initialize blockchain with validators and other info from TendermintCore
func (app *MentaApp) InitChain(req abci.RequestInitChain) abci.ResponseInitChain {
	// TODO: add Validators to state
	if app.onGenesisHandler != nil {
		ctx := sdk.NewContext(app.deliverCache, nil)
		app.onGenesisHandler(ctx, req)
	}
	return abci.ResponseInitChain{}
}

func (app *MentaApp) BeginBlock(req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	if app.onBeginBlockHandler != nil {
		ctx := sdk.NewContext(app.deliverCache, nil)
		app.onBeginBlockHandler(ctx, req)
	}
	return abci.ResponseBeginBlock{}
}

func (app *MentaApp) DeliverTx(raw []byte) abci.ResponseDeliverTx {
	tx, err := sdk.TransactionFromBytes(raw)
	if err != nil {
		e := sdk.ErrorBadTx()
		return abci.ResponseDeliverTx{Code: e.Code, Log: e.Log}
	}

	// Check for the handler
	handler := app.handlers[tx.Call]
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

// Signals the end of a block, returns changes to the validator set
func (app *MentaApp) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	if app.onEndBlockHandler != nil {
		ctx := sdk.NewContext(app.deliverCache, nil)
		return app.onEndBlockHandler(ctx, req)
	}
	return abci.ResponseEndBlock{}
}

// Commit to state, refresh cache
func (app *MentaApp) Commit() abci.ResponseCommit {
	app.deliverCache.ApplyToState()

	commitresults := app.state.Commit()

	app.deliverCache = app.state.RefreshCache().(*store.KVCache)
	app.checkCache = app.state.RefreshCache().(*store.KVCache)

	return abci.ResponseCommit{Data: commitresults.Hash}
}
