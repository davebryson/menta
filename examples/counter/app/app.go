package app

// Classis example of the counter application
// implemented with Menta

import (
	"encoding/binary"
	fmt "fmt"

	menta "github.com/davebryson/menta/app"
	sdk "github.com/davebryson/menta/types"
)

const (
	// RouteName is the tx route
	serviceName = "counter"
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

// CounterExampleService implementation for counter
type CounterExampleService struct{}

// Route returns the registered name of the service
func (srv CounterExampleService) Route() string { return serviceName }

// Init set the initial state
func (srv CounterExampleService) Init(ctx sdk.TxContext) {
	ctx.Insert(stateKey, encodeCount(0))
}

// Execute runs the core logic
func (srv CounterExampleService) Execute(ctx sdk.TxContext) sdk.Result {
	ctx.Insert(stateKey, ctx.Tx.Msg)
	return sdk.Result{
		Log: "ok",
	}
}

// Query returns the current committed state
func (srv CounterExampleService) Query(key []byte, ctx sdk.QueryContext) sdk.Result {
	val, err := ctx.Get(key)
	if err != nil {
		return sdk.ResultError(1, err.Error())
	}
	return sdk.Result{
		Code: 0,
		Data: val,
	}
}

// ValidateCounterTx valdates new transactions
func ValidateCounterTx(ctx sdk.TxContext) sdk.Result {
	// Decode the incoming msg in the Tx
	msgVal := decodeCount(ctx.Tx.Msg)

	// Decode the state
	val, err := ctx.Get(stateKey)
	if err != nil {
		return sdk.ResultError(2, "expected count")
	}
	stateCount := decodeCount(val)

	// msg should match the expected next state
	expected := stateCount + uint32(1)
	if msgVal != expected {
		return sdk.ResultError(2, fmt.Sprintf("Error: bad count expected %d", expected))
	}

	// Increment the state so other checks are correct
	ctx.Insert(stateKey, encodeCount(msgVal))

	return sdk.Result{
		Log: "ok",
	}
}

// This is the counter application
func createApp() *menta.MentaApp {
	// runs tendermint init if needed
	menta.InitTendermint(HomeDir)
	// setup the app
	app := menta.NewApp("counter-example", HomeDir)
	// add the check tx handler
	app.ValidateTxHandler(ValidateCounterTx)
	// add the service
	app.AddService(CounterExampleService{})

	return app
}

// RunApp sets up the menta application and starts the node
func RunApp() {
	app := createApp()
	app.Run()
}
