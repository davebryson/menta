package revamp

import (
	"testing"

	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
	amino "github.com/tendermint/go-amino"
)

type CreateMsg struct {
	Name string
}

func (cm CreateMsg) Execute(sender []byte, store sdk.KVStore) sdk.Result {
	store.Set([]byte("name"), []byte(cm.Name))

	return sdk.Result{
		Code: 0,
		Log:  "Got it",
	}
}

// This approach (compared to a handler) allows each service to also
// have some transient state.
type MyService struct{}

func (ms *MyService) Route() string {
	return "myservice"
}
func (ms *MyService) RegisterMessages(cdc *amino.Codec) {
	cdc.RegisterConcrete(&CreateMsg{}, "myservice/cm", nil)
}
func (ms *MyService) OnGenesis(store sdk.KVStore)         {}
func (ms *MyService) Query(key []byte, store sdk.KVStore) {}
func (ms *MyService) Execute(sender []byte, msg Message, store sdk.KVStore) sdk.Result {
	return msg.Execute(sender, store)
}

func MakeTx(name string, cdc *amino.Codec) []byte {
	tx := SignedTransaction{
		Route:     "example1",
		Sender:    []byte("dave"),
		Msg:       CreateMsg{Name: name},
		Nonce:     []byte("1"),
		Signature: []byte("sig"),
	}
	return cdc.MustMarshalBinaryLengthPrefixed(tx)
}

func TestServices(t *testing.T) {
	assert := assert.New(t)
	simapp := NewSimApp()
	simapp.AddService(&MyService{})

	tx1 := MakeTx("dave", simapp.codec)

	result := simapp.Run(tx1)
	assert.Equal(uint32(0), result.Code)
	assert.Equal("Got it", result.Log)
	v := simapp.dcache.Get([]byte("name"))
	assert.NotNil(v)
	assert.Equal([]byte("dave"), v)

	tx1 = MakeTx("bob", simapp.codec)
	result = simapp.Run(tx1)
	assert.Equal(uint32(0), result.Code)
	assert.Equal("Got it", result.Log)
	v = simapp.dcache.Get([]byte("name"))
	assert.NotNil(v)
	assert.Equal([]byte("bob"), v)

}
