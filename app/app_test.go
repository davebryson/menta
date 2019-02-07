package app

import (
	"encoding/binary"
	"testing"

	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func bigE(v uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}

func makeTx() ([]byte, error) {
	t := &sdk.Transaction{Call: "counter_test"}
	return t.Bytes()
}

// Test to check all callbacks and handler hooks
func TestAppCallbacks(t *testing.T) {
	assert := assert.New(t)

	// Setup the app
	app := NewMockApp() // inmemory tree

	// Set up initial chain state
	app.OnInitialStart(func(ctx sdk.Context, req abci.RequestInitChain) (resp abci.ResponseInitChain) {
		ctx.Db.Set([]byte("count"), bigE(1))
		return
	})
	// Add a validator
	app.OnValidateTx(func(ctx sdk.Context) sdk.Result {
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
	// Increments the count and updates state
	app.OnTx("counter_test", func(ctx sdk.Context) sdk.Result {
		count := binary.BigEndian.Uint32(ctx.Db.Get([]byte("count")))
		count += uint32(1)
		ctx.Db.Set([]byte("count"), bigE(count))
		return sdk.Result{
			Log: "increment",
		}
	})
	// Add an EndBlock handler
	app.OnEndBlock(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return abci.ResponseEndBlock{
			Tags: sdk.Tags{
				sdk.Tag{Key: []byte("end"), Value: []byte("av")},
			},
		}
	})

	// Simulate running it

	// Call InitChain
	icresult := app.InitChain(abci.RequestInitChain{})
	assert.Equal(icresult, abci.ResponseInitChain{})
	// Commit here so we have something for Info
	c1 := app.Commit()

	// Call Info
	result := app.Info(abci.RequestInfo{})
	assert.Equal("mockapp", result.GetData())
	// Should be 1 because we committed
	assert.Equal(int64(1), result.GetLastBlockHeight())
	hash1 := result.GetLastBlockAppHash()
	assert.NotNil(hash1)
	// Should == the first commit hash
	assert.Equal(c1.Data, hash1)

	// Call Query & Check the state
	respQ := app.Query(abci.RequestQuery{Path: "/key", Data: []byte("count")})
	assert.Equal(uint32(0), respQ.Code)
	assert.Equal(uint32(1), binary.BigEndian.Uint32(respQ.GetValue()))

	// Run validate
	tx, err := makeTx()
	assert.Nil(err)
	chtx := app.CheckTx(tx)
	assert.Equal(abci.ResponseCheckTx{Log: "ok", Code: 0}, chtx)

	bb := app.BeginBlock(abci.RequestBeginBlock{})
	assert.Equal(1, len(bb.Tags))
	assert.Equal([]byte("begin"), bb.Tags[0].Key)

	// Run Tx handler
	dtx := app.DeliverTx(tx)
	assert.Equal(abci.ResponseDeliverTx{Log: "increment", Code: 0}, dtx)

	eb := app.EndBlock(abci.RequestEndBlock{})
	assert.Equal(1, len(eb.Tags))
	assert.Equal([]byte("end"), eb.Tags[0].Key)

	// Commit the new state
	commit := app.Commit()
	assert.NotNil(commit.Data)
	assert.NotEqual(c1.Data, commit.Data)

	// Now state should == 1
	respQ = app.Query(abci.RequestQuery{Path: "/key", Data: []byte("count")})
	assert.Equal(uint32(1), binary.BigEndian.Uint32(respQ.GetValue()))
}
