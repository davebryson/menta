package main

import (
	"fmt"
	"os"

	menta "github.com/davebryson/menta/app"
	"github.com/davebryson/menta/codec"
	"github.com/davebryson/menta/store"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var countKey = []byte("countvalue")

// StoreValue is an example struct for storing the value and returning it from a query
type StoreValue struct {
	Current int `json:"current"`
}

// CounterMsg is the tx msg
type CounterMsg struct {
	Add int `json:"add"`
}

// Route ..
func (msg CounterMsg) Route() string { return "counterapp" }

// Type ..
func (msg CounterMsg) Type() string { return "" }

// ValidateBasic ..
func (msg CounterMsg) ValidateBasic() error {
	if msg.Add < 0 {
		return fmt.Errorf("Must be >= zero")
	}
	return nil
}

func main() {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdk.RegisterStandardTypes(cdc)
	cdc.RegisterConcrete(&CounterMsg{}, "menta/counterMsg", nil)

	// Runs tendermint init
	homedir := os.ExpandEnv(fmt.Sprintf("$HOME/%s", ".counterapp"))
	menta.InitTendermint(homedir)

	// Setup the app
	app := menta.NewApp("counterapp", homedir, sdk.DefaultJSONTxDecoder(cdc))

	// Setup initial state initial
	app.OnInitialStart(func(ctx sdk.Context, req abci.RequestInitChain) (res abci.ResponseInitChain) {
		// Set the initial count value to 0
		ctx.Db.Set(countKey, cdc.MustMarshalBinaryBare(StoreValue{Current: 0}))
		return
	})

	// Route tx handler to add to the last state
	app.Route("counter", func(ctx sdk.Context) sdk.Result {
		switch msg := ctx.Tx.GetMsg().(type) {
		case CounterMsg:
			// decode from storage
			var v StoreValue
			if err := cdc.UnmarshalBinaryBare(ctx.Db.Get(countKey), &v); err != nil {
				return sdk.ErrorBadTx()
			}

			// Add to the state
			v.Current += msg.Add
			// Save it
			ctx.Db.Set(countKey, cdc.MustMarshalBinaryBare(v))
			return sdk.Result{}

		default:
			return sdk.ResultError(10, "Unknown message")
		}
	})

	// RouteQuery to see the latest count from committed storage
	app.RouteQuery("current/count", func(key []byte, ctx store.QueryContext) abci.ResponseQuery {
		res := abci.ResponseQuery{}
		res.Code = 10 // fail

		data := ctx.Get(countKey)
		if data == nil {
			res.Log = "Can't find count in storage"
			return res
		}

		jsonbits, err := cdc.MarshalJSON(data)
		if err != nil {
			res.Log = err.Error()
			return res
		}

		res.Code = 0
		res.Value = jsonbits
		return res
	})

	// run with the app embedded in Tendermint
	app.Run()
}
