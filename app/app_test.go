package app

import (
	"encoding/binary"
	"testing"

	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	stateKey    = []byte("countkey")
	serviceName = "counter_test"
)

var _ sdk.Service = (*CounterService)(nil)

func encodeCount(val uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	return buf
}

func decodeCount(bits []byte) uint32 {
	return binary.BigEndian.Uint32(bits)
}

func makeTx(val uint32) ([]byte, error) {
	encoded := encodeCount(val)
	t := &sdk.SignedTransaction{Service: serviceName, Msg: encoded}
	return sdk.EncodeTx(t)
}

// Service implementation

type CounterService struct{}

func (srv CounterService) Route() string { return serviceName }
func (srv CounterService) Init(ctx sdk.TxContext) {
	ctx.Insert(stateKey, encodeCount(0))
}
func (srv CounterService) Execute(ctx sdk.TxContext) sdk.Result {
	ctx.Insert(stateKey, ctx.Tx.Msg)
	return sdk.Result{
		Log: "ok",
	}
}
func (srv CounterService) Query(key []byte, ctx sdk.QueryContext) sdk.Result {
	val, err := ctx.Get(key)
	if err != nil {
		return sdk.ResultError(1, err.Error())
	}
	return sdk.Result{
		Code: 0,
		Data: val,
	}
}

func ValidateTx(ctx sdk.TxContext) sdk.Result {
	// Decode the incoming msg in the Tx
	msgVal := decodeCount(ctx.Tx.Msg)
	// Decode the state
	val, err := ctx.Get(stateKey)
	if err != nil {
		return sdk.ResultError(2, "expected count")
	}
	stateCount := decodeCount(val)

	// msg should match the expected next state
	expected := stateCount + uint32(1)
	if msgVal != expected {
		return sdk.ResultError(2, "bad count")
	}

	// Increment the state for other checks
	ctx.Insert(stateKey, encodeCount(msgVal))

	return sdk.Result{
		Log: "ok",
	}
}

func createApp() *MentaApp {
	app := NewMockApp() // inmemory tree
	app.ValidateTxHandler(ValidateTx)
	app.AddService(CounterService{})
	return app
}

// Test to check all callbacks and handler hooks
func TestAppCallbacks(t *testing.T) {
	assert := assert.New(t)
	app := createApp()

	// --- Simulate running it ---

	// Call InitChain
	icresult := app.InitChain(abci.RequestInitChain{})
	assert.Equal(icresult, abci.ResponseInitChain{})
	// Commit here so we have something for Info
	c1 := app.Commit()

	// Call Info
	result := app.Info(abci.RequestInfo{})
	assert.Equal("mockapp", result.GetData())
	// block height should be 1 because we committed
	assert.Equal(int64(1), result.GetLastBlockHeight())
	hash1 := result.GetLastBlockAppHash()
	assert.NotNil(hash1)
	// Should == the first commit hash
	assert.Equal(c1.Data, hash1)

	// Call Query & check the initial state
	respQ := app.Query(abci.RequestQuery{Path: serviceName, Data: stateKey})
	assert.Equal(uint32(0), respQ.Code)
	assert.Equal(uint32(0), decodeCount(respQ.GetValue()))

	// Run check
	// Ok
	tx, err := makeTx(1)
	assert.Nil(err)
	chtx := app.CheckTx(abci.RequestCheckTx{Tx: tx})
	assert.Equal(abci.ResponseCheckTx{Log: "ok", Code: 0}, chtx)

	// Run check
	// Bad count - should be 2
	badtx, _ := makeTx(4)
	chtx = app.CheckTx(abci.RequestCheckTx{Tx: badtx})
	assert.Equal(abci.ResponseCheckTx{Log: "bad count", Code: 2}, chtx)

	// Run check
	// Should be good
	tx1, err := makeTx(2)
	chtx = app.CheckTx(abci.RequestCheckTx{Tx: tx1})
	assert.Equal(abci.ResponseCheckTx{Log: "ok", Code: 0}, chtx)

	// Check cache should return 0 as it only reads committed state
	respQ = app.Query(abci.RequestQuery{Path: serviceName, Data: stateKey})
	assert.Equal(uint32(0), respQ.Code)
	assert.Equal(uint32(0), decodeCount(respQ.GetValue()))

	// Run Deliver handlers
	dtx := app.DeliverTx(abci.RequestDeliverTx{Tx: tx})
	assert.Equal(abci.ResponseDeliverTx{Log: "ok", Code: 0}, dtx)

	dtx = app.DeliverTx(abci.RequestDeliverTx{Tx: tx1})
	assert.Equal(abci.ResponseDeliverTx{Log: "ok", Code: 0}, dtx)

	// Commit the new state to storage
	commit := app.Commit()

	assert.NotNil(commit.Data)
	// Should be a new apphash
	assert.NotEqual(c1.Data, commit.Data)

	// Now committed state should == 2
	respQ = app.Query(abci.RequestQuery{Path: serviceName, Data: stateKey})
	assert.Equal(uint32(2), decodeCount(respQ.GetValue()))
}
