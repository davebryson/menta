package app

import (
	"encoding/binary"
	fmt "fmt"

	menta "github.com/davebryson/menta/app"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// RouteName is the tx route
	routeName = "counter"
	// QueryRoute is the route name for the query handler
	queryRoute = "showstate"
	// HomeDir is the name of the directory where all the tendermint data is stored for the node
	HomeDir = "counterdata"
)

// Key used to reference our count data in storage
var stateKey = []byte("counterStateKey")

// Encode the value
func encodeCount(val uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	return buf
}

// decode the value
func decodeCount(bits []byte) uint32 {
	return binary.BigEndian.Uint32(bits)
}

// This is the counter application
func createApp() *menta.MentaApp {
	// runs tendermint init if needed
	menta.InitTendermint(HomeDir)

	// setup the app
	app := menta.NewApp("counter-example", HomeDir)

	// Set the initial state to 0. This is only ran once the first time the app is started
	app.OnInitChain(func(store sdk.RWStore, req abci.RequestInitChain) (resp abci.ResponseInitChain) {
		store.Set(stateKey, encodeCount(0))
		return
	})

	// Add the tx validator
	app.OnValidateTx(func(store sdk.RWStore, tx *sdk.Tx) sdk.Result {
		// Decode the incoming msg in the Tx
		msgVal := decodeCount(tx.Msg)

		v, err := store.Get(stateKey)
		if err != nil {
			return sdk.ResultError(2, err.Error())
		}

		// Decode the state
		stateCount := decodeCount(v)

		// msg should match the expected next state
		expected := stateCount + uint32(1)
		if msgVal != expected {
			return sdk.ResultError(2, fmt.Sprintf("Error: bad count expected %d", expected))
		}

		// Increment the state so other checks are correct
		store.Set(stateKey, encodeCount(msgVal))

		return sdk.Result{
			Log: "ok",
		}
	})

	// Add a Tx handler to update state
	// Sets state to the value of tx.msg.  This is ok as we've already validated the tx
	// in checkTx
	app.OnTx(routeName, func(store sdk.RWStore, tx *sdk.Tx) sdk.Result {
		store.Set(stateKey, tx.Msg)
		return sdk.Result{
			Log: "increment",
		}
	})

	// Handle queries for the current committed state
	app.OnQuery(queryRoute, func(store sdk.StoreReader, key []byte) ([]byte, error) {
		return store.Get(key)
	})

	return app
}

// RunApp sets up the menta application and starts the node
func RunApp() {
	app := createApp()
	app.Run()
}
