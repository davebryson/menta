# MENTA
A *simple* framework for creating Tendermint blockchain applications.  

But what about [Cosmos-SDK](https://github.com/cosmos/cosmos-sdk) ?.  Good question, if you're building a public application or planning to deploy to Cosmos, you should definitely use it.  

This framework is for:
* rapid prototyping and small pilot projects  
* folks looking to build Tendermint apps not destined for the Cosmos, or 
* folks wanting to learn *how* a Tendermint ABCI works

Of course, you can always start here and port to the Cosmos SDK later. That's the magic of ABCI

**Current supported Tendermint version: 0.29.0**

Get it: `dep ensure -add github.com/davebryson/menta`

Get started:
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

	// Set up the initial state for the new app - (this is the abci.InitChain())
	app.OnInitialStart(func(ctx sdk.Context, req abci.RequestInitChain) {
		ctx.Db.Set([]byte("count"), writeNumber(0))
	})

	// Add 1 to many transaction callbacks (this is the abci.DeliverTx)
	app.OnTx("counter", func(ctx sdk.Context) sdk.Result {
		count := binary.BigEndian.Uint32(ctx.Db.Get([]byte("count")))
		count += uint32(1)
		ctx.Db.Set([]byte("count"), writeNumber(count))
		return sdk.Result{}
	})

	// Run with the app - embedded with Tendermint
	app.Run()
}
```
