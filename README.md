# MENTA
A simple golang framework for creating Tendermint blockchain applications.

Get it: `dep ensure -add github.com/davebryson/menta`

Here's an example of the ubiquitous 'counter' application:
```golang
package main

import (
	"encoding/binary"

	menta "github.com/davebryson/menta/app"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/abci/types"
)

// Formatting stuff...
func writeNumber(v uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}

func main() {

	// Runs tendermint init - "" default to ~/.menta
	menta.InitTendermint("")

	// Setup the application. "counterapp is the name-o"
	app := menta.NewApp("counterapp", "")

	// Set up the initial - (this is the abci.InitChain())
	app.OnGenesis(func(ctx sdk.Context, req abci.RequestInitChain) {
		ctx.Db.Set([]byte("count"), writeNumber(0))
	})

	// Add 1 to many transaction callbacks (this is the abci.DeliverTx)
	app.OnTx("counter", func(ctx sdk.Context) sdk.Result {
		count := binary.BigEndian.Uint32(ctx.Db.Get([]byte("count")))
		count += uint32(1)
		ctx.Db.Set([]byte("count"), writeNumber(count))
		return sdk.Result{}
	})

	// Run with the app - embedded in Tendermint
	app.Run()

}
```

Menta also makes local development and testing a snap!

Here's how to test the app without running a node:
```golang
// in some test function...
func TestTesterBasics(t *testing.T) {
	assert := assert.New(t)

        // Expected tree hash after the transactions
	expectedTreeHash := "cd5c39b2ff82a4fe914b095daef7becc348709c5"

	// Setup app - not the use of 'NewTestApp'
	app := menta.NewTestApp()

        // Same as the example above
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

	// Run a query - initial state should be zero
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
```

More to come.  Of course if you're looking for a more full-featured SDK, I'd recommend [Cosmos-SDK](https://github.com/cosmos/cosmos-sdk)

TODO: Upgrade to latest ABCI 12
