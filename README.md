# MENTA
A *simple* framework for creating Tendermint blockchain applications.  

Why not just use the [Cosmos-SDK](https://github.com/cosmos/cosmos-sdk) ?.  Good question. If you're building a public application or planning to deploy to Cosmos, you should definitely use it!  

This framework is for:
* rapid prototyping and small pilot projects  
* folks looking to build Tendermint apps that may not be destined for the Cosmos, or 
* folks wanting to learn *how* a Tendermint ABCI works through a relatively small code base.

Menta intentionally mimics, adapts, and uses code from the Cosmos SDK.  As I use this framework as a way to help me stay in tune with what's happening under the covers with Tendermint and the Cosmos SDK.

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
