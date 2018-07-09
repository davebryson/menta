package main

import (
	"encoding/binary"

	menta "github.com/davebryson/menta/app"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func writeNumber(v uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}

func main() {

	// runs tendermint init - "" default to ~/.menta
	menta.InitTendermint("")

	// setup the app
	app := menta.NewApp("counterapp", "")

	// initial state callback
	app.OnGenesis(func(ctx sdk.Context, req abci.RequestInitChain) {
		ctx.Db.Set([]byte("count"), writeNumber(0))
	})

	// tx callback to increment the count
	app.OnTx("counter", func(ctx sdk.Context) sdk.Result {
		count := binary.BigEndian.Uint32(ctx.Db.Get([]byte("count")))
		count += uint32(1)
		ctx.Db.Set([]byte("count"), writeNumber(count))
		return sdk.Result{}
	})

	// run with the app embedded in Tendermint
	app.Run()

}
