package test

import (
	"encoding/binary"
	"fmt"
	"testing"

	menta "github.com/davebryson/menta/app"
	"github.com/davebryson/menta/tools"
	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	//ex "github.com/davebryson/menta/examples"
)

func writeNumber(v uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}

func counterInitState(ctx sdk.Context, req abci.RequestInitChain) {
	ctx.Db.Set([]byte("count"), writeNumber(0)) // limited to 255
}

func counterTx(ctx sdk.Context) sdk.Result {
	count := binary.BigEndian.Uint32(ctx.Db.Get([]byte("count")))
	count += uint32(1)
	ctx.Db.Set([]byte("count"), writeNumber(count))
	return sdk.Result{}
}

func makeTx() ([]byte, error) {
	t := &sdk.Transaction{Call: "counter"}
	return t.Bytes()
}

func TestApp(t *testing.T) {
	assert := assert.New(t)

	// Setup the chain
	app := menta.NewTestApp() // memdb
	app.OnTx("counter", counterTx)
	app.OnGenesis(counterInitState)
	// Setup initial chain - must be done first
	app.InitChain(abci.RequestInitChain{})
	app.Commit()

	// Info
	result := app.Info(abci.RequestInfo{})
	assert.Equal("testapp", result.GetData())
	assert.Equal(int64(1), result.GetLastBlockHeight())
	assert.NotNil(result.GetLastBlockAppHash())

	// Query
	respQ := app.Query(abci.RequestQuery{Path: "/key", Data: []byte("count")})
	assert.Equal(uint32(0), binary.BigEndian.Uint32(respQ.GetValue()))

	// Do 10 txs
	for i := 0; i < 10; i++ {
		tx1, err := makeTx()
		assert.Nil(err)
		rD := app.DeliverTx(tx1)
		assert.Equal(uint32(0), rD.GetCode())
	}

	// Before commit should still == 0
	respQ = app.Query(abci.RequestQuery{Path: "/key", Data: []byte("count")})
	assert.Equal(uint32(0), binary.BigEndian.Uint32(respQ.GetValue()))

	app.Commit()

	// Now it is the last state
	respQ = app.Query(abci.RequestQuery{Path: "/key", Data: []byte("count")})
	assert.Equal(uint32(10), binary.BigEndian.Uint32(respQ.GetValue()))
}

func TestTesterBasics(t *testing.T) {
	assert := assert.New(t)

	expectedTreeHash := "cd5c39b2ff82a4fe914b095daef7becc348709c5"

	// Setup app
	app := menta.NewTestApp() // memdb
	app.OnGenesis(func(ctx sdk.Context, req abci.RequestInitChain) {
		ctx.Db.Set([]byte("count"), writeNumber(0))
	})
	app.OnTx("counter", func(ctx sdk.Context) sdk.Result {
		count := binary.BigEndian.Uint32(ctx.Db.Get([]byte("count")))
		count += uint32(1)
		ctx.Db.Set([]byte("count"), writeNumber(count))
		return sdk.Result{}
	})

	// Setup tester
	tapp := tools.NewTester(app)

	// Initial state should be zero
	_, val := tapp.QueryByKey("count")
	assert.Equal(uint32(0), binary.BigEndian.Uint32(val))

	// Do some Txs
	for i := 0; i < 10; i++ {
		tx1, err := makeTx()
		assert.Nil(err)
		tapp.SendTx(tx1)
	}

	// Make Block
	c1, numTx := tapp.MakeBlock()
	hash := fmt.Sprintf("%x", c1.GetHash())

	assert.Equal(10, numTx)
	assert.Equal(int64(1), c1.GetVersion())
	assert.Equal(expectedTreeHash, hash)

	_, val = tapp.QueryByKey("count")
	assert.Equal(uint32(10), binary.BigEndian.Uint32(val))
}
