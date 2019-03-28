package main

import (
	"encoding/binary"
	"fmt"

	menta "github.com/davebryson/menta/app"
	"github.com/davebryson/menta/codec"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func writeNumber(v uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}

var countKey = []byte("countvalue")

// CounterMsg is the tx msg
type CounterMsg struct {
	Value int
}

// Route ..
func (msg CounterMsg) Route() string { return "counterapp" }

// Type ..
func (msg CounterMsg) Type() string { return "" }

// ValidateBasic ..
func (msg CounterMsg) ValidateBasic() error {
	if msg.Value < 0 {
		return fmt.Errorf("Must be >= zero")
	}
	return nil
}

func main() {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdk.RegisterStandardTypes(cdc)
	cdc.RegisterConcrete(&CounterMsg{}, "menta/counterMsg", nil)

	// Runs tendermint init - "" default to ~/.menta
	menta.InitTendermint("")

	// Setup the app
	app := menta.NewApp("counterapp", "", sdk.DefaultJSONTxDecoder(cdc))

	// Setup initial state initial
	app.OnInitialStart(func(ctx sdk.Context, req abci.RequestInitChain) (res abci.ResponseInitChain) {
		// Set the initial count value to 0
		ctx.Db.Set(countKey, cdc.MustMarshalBinaryBare(0))
		return
	})

	// Route that adds to the last state
	app.Route("counter", func(ctx sdk.Context) sdk.Result {
		switch msg := ctx.Tx.GetMsg().(type) {
		case CounterMsg:
			// decode from storage
			var v int
			if err := cdc.UnmarshalBinaryBare(ctx.Db.Get(countKey), &v); err != nil {
				return sdk.ErrorBadTx()
			}

			// Add to the state
			v += msg.Value
			// Save it
			ctx.Db.Set(countKey, cdc.MustMarshalBinaryBare(v))
			return sdk.Result{}

		default:
			return sdk.ResultError(10, "Unknown message")
		}
	})

	// run with the app embedded in Tendermint
	app.Run()
}
