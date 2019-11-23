package app

import (
	"encoding/binary"
	"testing"

	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	stateKey = []byte("countkey")
)

const routeName = "counter_test"

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
	t := &sdk.Tx{Route: routeName, Msg: encoded}
	return sdk.EncodeTx(t)
}

func createApp() *MentaApp {
	app := NewMockApp() // inmemory tree

	// Set up initial chain state
	app.OnInitChain(func(store sdk.RWStore, req abci.RequestInitChain) (resp abci.ResponseInitChain) {
		store.Set(stateKey, encodeCount(0))
		return
	})
	// Add the tx validator
	app.OnValidateTx(func(store sdk.RWStore, tx *sdk.Tx) sdk.Result {
		// Decode the incoming msg in the Tx
		msgVal := decodeCount(tx.Msg)
		// Decode the state
		v, err := store.Get(stateKey)
		if err != nil {
			return sdk.ResultError(2, err.Error())
		}
		stateCount := decodeCount(v)
		// msg should match the expected next state
		expected := stateCount + uint32(1)
		if msgVal != expected {
			return sdk.ResultError(2, "bad count")
		}

		// Increment the state for other checks
		store.Set(stateKey, encodeCount(msgVal))

		return sdk.Result{
			Log: "ok",
		}
	})
	// Add a BeginBlock handler
	app.OnBeginBlock(func(store sdk.RWStore, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		return abci.ResponseBeginBlock{
			//Tags: sdk.Tags{
			//	sdk.Tag{Key: []byte("begin"), Value: []byte("av")},
			//},
		}
	})
	// Add a Tx processor with 'counter_test' route
	// Increments the count from the msg and updates state
	app.OnTx(routeName, func(store sdk.RWStore, tx *sdk.Tx) sdk.Result {
		store.Set(stateKey, tx.Msg)
		return sdk.Result{
			Log: "increment",
		}
	})

	app.OnQuery("/key", func(store sdk.StoreReader, key []byte) ([]byte, error) {
		return store.Get(key)
	})

	// Add an EndBlock handler
	app.OnEndBlock(func(store sdk.RWStore, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return abci.ResponseEndBlock{
			//Tags: sdk.Tags{
			//	sdk.Tag{Key: []byte("end"), Value: []byte("av")},
			//},
		}
	})
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

	// Call Query & Check the state
	respQ := app.Query(abci.RequestQuery{Path: "/key", Data: stateKey})
	assert.Equal(uint32(0), respQ.Code)
	assert.Equal(uint32(0), decodeCount(respQ.GetValue()))

	// Run validate
	// Ok
	tx, err := makeTx(1)
	assert.Nil(err)
	chtx := app.CheckTx(abci.RequestCheckTx{Tx: tx})
	assert.Equal(abci.ResponseCheckTx{Log: "ok", Code: 0}, chtx)

	// Bad
	badtx, _ := makeTx(4)
	chtx = app.CheckTx(abci.RequestCheckTx{Tx: badtx})
	assert.Equal(abci.ResponseCheckTx{Log: "bad count", Code: 2}, chtx)

	// Ok
	tx1, err := makeTx(2)
	chtx = app.CheckTx(abci.RequestCheckTx{Tx: tx1})
	assert.Equal(abci.ResponseCheckTx{Log: "ok", Code: 0}, chtx)

	// No committed state yet. So it should still be 0
	respQ = app.Query(abci.RequestQuery{Path: "/key", Data: stateKey})
	assert.Equal(uint32(0), respQ.Code)
	assert.Equal(uint32(0), decodeCount(respQ.GetValue()))

	// --- Process a block
	//bb := app.BeginBlock(abci.RequestBeginBlock{})
	//assert.Equal(1, len(bb.Tags))
	//assert.Equal([]byte("begin"), bb.Tags[0].Key)

	// Run Tx handlers
	dtx := app.DeliverTx(abci.RequestDeliverTx{Tx: tx})
	assert.Equal(abci.ResponseDeliverTx{Log: "increment", Code: 0}, dtx)
	dtx = app.DeliverTx(abci.RequestDeliverTx{Tx: tx1})
	assert.Equal(abci.ResponseDeliverTx{Log: "increment", Code: 0}, dtx)

	//eb := app.EndBlock(abci.RequestEndBlock{})
	//assert.Equal(1, len(eb.Tags))
	//assert.Equal([]byte("end"), eb.Tags[0].Key)

	// Commit the new state to storage
	commit := app.Commit()
	assert.NotNil(commit.Data)
	// Should be a new apphash
	assert.NotEqual(c1.Data, commit.Data)

	// Now committed state should == 2
	respQ = app.Query(abci.RequestQuery{Path: "/key", Data: stateKey})
	assert.Equal(uint32(2), decodeCount(respQ.GetValue()))
}
