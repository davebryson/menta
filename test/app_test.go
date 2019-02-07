package test

import (
	"encoding/binary"
	"testing"

	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func writeNumber(v uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}

func counterInitState(ctx sdk.Context, req abci.RequestInitChain) (resp abci.ResponseInitChain) {
	ctx.Db.Set([]byte("count"), writeNumber(0)) // limited to 255
	return
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

func TestTesterBasics(t *testing.T) {
	/*assert := assert.New(t)

	expectedTreeHash := "d25242c980689eafc4ead9fb60f64f7542a03c44632af735bd5e216bad3aa8da"

	// Setup app
	app := menta.NewMockApp() // memdb
	app.OnInitialStart(
		func(ctx sdk.Context, req abci.RequestInitChain) (res abci.ResponseInitChain) {
			ctx.Db.Set([]byte("count"), writeNumber(0))
			return
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
	assert.Equal(uint32(9), binary.BigEndian.Uint32(val))*/
}
