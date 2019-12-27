package counter

import (
	sdk "github.com/davebryson/menta/types"
	"github.com/gogo/protobuf/proto"
)

// ServiceName is just that...
const ServiceName = "counter_example"

var _ sdk.Service = (*Service)(nil)

// Service is a simple service to demonstrate
// the menta API.  It stores an counter for each tx.sender
type Service struct{}

// Name is a unique name used to register the service
func (srv Service) Name() string { return ServiceName }

// Initialize is called on the genesis block.  Not used
func (srv Service) Initialize(data []byte, store sdk.KVStore) {
}

// Execute runs the core logic for a state transition
func (srv Service) Execute(sender []byte, msgid uint32, message []byte, store sdk.KVStore) sdk.Result {
	// Decode the incoming msg in the Tx
	var msg Increment
	err := proto.Unmarshal(message, &msg)
	if err != nil {
		return sdk.ErrorBadTx()
	}

	schema := NewSchema(store)
	return schema.IncrementCount(sender, msg)
}

// Query committed state for the given used. Key is the public key bytes
func (srv Service) Query(key []byte, store sdk.KVStore) sdk.Result {
	schema := NewSchema(store)
	return schema.Query(key)
}

type Schema struct {
	store sdk.NamedStore
}

func NewSchema(store sdk.KVStore) Schema {
	return Schema{
		store: sdk.NewNamedStore(ServiceName, store),
	}
}

func (schema Schema) IncrementCount(sender []byte, msg Increment) sdk.Result {
	storeVal, err := schema.store.Get(sender)
	if err != nil {
		// First tx
		msg, err := NewCounter(1).Encode()
		if err != nil {
			return sdk.ResultError(1, "problem encoding new count value")
		}
		schema.store.Put(sender, msg)
		return sdk.Result{}
	}

	// Decode the current value in the store
	stateCount, err := DecodeCount(storeVal)
	if err != nil {
		return sdk.ResultError(1, "problem decoding stored value")
	}

	// 'tx.msg' should match the expected next state
	if !stateCount.ValidNextValue(msg.Value) {
		return sdk.ResultError(2, "bad count")
	}

	// Increment the count and update storage
	stateCount.Inc()
	newcount, err := stateCount.Encode()
	if err != nil {
		return sdk.ResultError(1, "problem encoding new count value")
	}

	// It's good, save it
	schema.store.Put(sender, newcount)

	return sdk.Result{
		Log: "ok",
	}
}

func (schema Schema) Query(key []byte) sdk.Result {
	val, err := schema.store.Query(key)
	if err != nil {
		return sdk.ResultError(1, err.Error())
	}
	return sdk.Result{
		Code: 0,
		Data: val,
	}
}

// --- Augment proto types ---

// Encode the Increment message
func (inc *Increment) Encode() ([]byte, error) {
	return proto.Marshal(inc)
}

// NewCounter set to 'val'
func NewCounter(val uint32) *CountValue {
	return &CountValue{
		Current: val,
	}
}

// Encode the Counter
func (count *CountValue) Encode() ([]byte, error) {
	return proto.Marshal(count)
}

// Inc - increments the counter by 1
func (count *CountValue) Inc() {
	count.Current++
}

// ValidNextValue check if the proposed count is correct
func (count *CountValue) ValidNextValue(proposed uint32) bool {
	return (count.Current + uint32(1)) == proposed
}

// DecodeCount bytes => Countvalue
func DecodeCount(raw []byte) (*CountValue, error) {
	var c CountValue
	err := proto.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
