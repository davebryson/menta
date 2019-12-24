package counter

import (
	sdk "github.com/davebryson/menta/types"
	"github.com/gogo/protobuf/proto"
)

const ServiceName = "counter_example"

var _ sdk.Service = (*CounterService)(nil)

type CounterService struct{}

func (srv CounterService) Name() string { return ServiceName }
func (srv CounterService) Initialize(data []byte, store sdk.PrefixedRWStore) {
}
func (srv CounterService) Execute(ctx sdk.TxContext) sdk.Result {
	// Decode the incoming msg in the Tx
	var msg Increment
	err := proto.Unmarshal(ctx.Message, &msg)
	if err != nil {
		return sdk.ErrorBadTx()
	}

	// If it's not found, and the message is 1. This is the first tx
	// for the user
	// Decode the current value in state
	storeVal, err := ctx.Store().Get(ctx.Sender)
	if err != nil {
		// First tx
		encodedCount, err := NewCounter(1).Encode()
		if err != nil {
			return sdk.ResultError(1, "problem encoding new count value")
		}
		ctx.Store().Insert(ctx.Sender, encodedCount)
		return sdk.Result{}
	}

	// Decode the current value in state
	stateCount, err := DecodeCount(storeVal)
	if err != nil {
		return sdk.ResultError(1, "problem decoding stored value")
	}

	// msg should match the expected next state
	if !stateCount.ValidNextValue(msg.Value) {
		return sdk.ResultError(2, "bad count")
	}

	// Inc
	stateCount.Inc()
	encodedCount, err := stateCount.Encode()
	if err != nil {
		return sdk.ResultError(1, "problem encoding new count value")
	}
	// It's good, save it
	ctx.Store().Insert(ctx.Sender, encodedCount)

	return sdk.Result{
		Log: "ok",
	}
}

func (srv CounterService) Query(key []byte, store sdk.PrefixedReadOnlyStore) sdk.Result {
	val, err := store.Get(key)
	if err != nil {
		return sdk.ResultError(1, err.Error())
	}
	return sdk.Result{
		Code: 0,
		Data: val,
	}
}

//

// --- Augment proto types ---

func (inc *Increment) Encode() ([]byte, error) {
	return proto.Marshal(inc)
}

func NewCounter(val uint32) *CountValue {
	return &CountValue{
		Current: val,
	}
}

func (count *CountValue) Encode() ([]byte, error) {
	return proto.Marshal(count)
}

func (count *CountValue) Inc() {
	count.Current++
}

func (count *CountValue) ValidNextValue(proposed uint32) bool {
	return (count.Current + uint32(1)) == proposed
}

func DecodeCount(raw []byte) (*CountValue, error) {
	var c CountValue
	err := proto.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
