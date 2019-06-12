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
	app.OnInitialStart(func(ctx sdk.Context, req abci.RequestInitChain) (resp abci.ResponseInitChain) {
		ctx.Db.Set(stateKey, encodeCount(0))
		return
	})
	// Add the tx validator
	app.OnValidateTx(func(ctx sdk.Context) sdk.Result {
		// Decode the incoming msg in the Tx
		msgVal := decodeCount(ctx.Tx.Msg)
		// Decode the state
		stateCount := decodeCount(ctx.Db.Get(stateKey))

		// msg should match the expected next state
		expected := stateCount + uint32(1)
		if msgVal != expected {
			return sdk.ResultError(2, "bad count")
		}

		// Increment the state for other checks
		ctx.Db.Set(stateKey, encodeCount(msgVal))

		return sdk.Result{
			Log: "ok",
		}
	})
	// Add a BeginBlock handler
	app.OnBeginBlock(func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		return abci.ResponseBeginBlock{
			Tags: sdk.Tags{
				sdk.Tag{Key: []byte("begin"), Value: []byte("av")},
			},
		}
	})
	// Add a Tx processor with 'counter_test' route
	// Increments the count from the msg and updates state
	app.OnTx(routeName, func(ctx sdk.Context) sdk.Result {
		ctx.Db.Set(stateKey, ctx.Tx.Msg)
		return sdk.Result{
			Log: "increment",
		}
	})

	app.OnQuery("/key", func(key []byte, ctx sdk.QueryContext) (resp abci.ResponseQuery) {
		resp.Value = ctx.Get(key)
		return
	})
	// Add an EndBlock handler
	app.OnEndBlock(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return abci.ResponseEndBlock{
			Tags: sdk.Tags{
				sdk.Tag{Key: []byte("end"), Value: []byte("av")},
			},
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
	chtx := app.CheckTx(tx)
	assert.Equal(abci.ResponseCheckTx{Log: "ok", Code: 0}, chtx)

	// Bad
	badtx, _ := makeTx(4)
	chtx = app.CheckTx(badtx)
	assert.Equal(abci.ResponseCheckTx{Log: "bad count", Code: 2}, chtx)

	// Ok
	tx1, err := makeTx(2)
	chtx = app.CheckTx(tx1)
	assert.Equal(abci.ResponseCheckTx{Log: "ok", Code: 0}, chtx)

	// No committed state yet. So it should still be 0
	respQ = app.Query(abci.RequestQuery{Path: "/key", Data: stateKey})
	assert.Equal(uint32(0), respQ.Code)
	assert.Equal(uint32(0), decodeCount(respQ.GetValue()))

	// --- Process a block
	bb := app.BeginBlock(abci.RequestBeginBlock{})
	assert.Equal(1, len(bb.Tags))
	assert.Equal([]byte("begin"), bb.Tags[0].Key)

	// Run Tx handlers
	dtx := app.DeliverTx(tx)
	assert.Equal(abci.ResponseDeliverTx{Log: "increment", Code: 0}, dtx)
	dtx = app.DeliverTx(tx1)
	assert.Equal(abci.ResponseDeliverTx{Log: "increment", Code: 0}, dtx)

	eb := app.EndBlock(abci.RequestEndBlock{})
	assert.Equal(1, len(eb.Tags))
	assert.Equal([]byte("end"), eb.Tags[0].Key)

	// Commit the new state to storage
	commit := app.Commit()
	assert.NotNil(commit.Data)
	// Should be a new apphash
	assert.NotEqual(c1.Data, commit.Data)

	// Now committed state should == 2
	respQ = app.Query(abci.RequestQuery{Path: "/key", Data: stateKey})
	assert.Equal(uint32(2), decodeCount(respQ.GetValue()))
}
