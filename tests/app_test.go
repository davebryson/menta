package test

import (
	"testing"

	auth "github.com/davebryson/menta/plugins/authenticate"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Test to check all callbacks and handler hooks
func TestAppCallbacks(t *testing.T) {
	assert := assert.New(t)

	bob := NewWallet()
	alice := NewWallet()

	app := createApp(bob, alice)

	// Simulate running it

	// Call InitChain
	icresult := app.InitChain(abci.RequestInitChain{})
	assert.Equal(icresult, abci.ResponseInitChain{})
	// Commit here so we have something for Info
	c1 := app.Commit()

	// Call Info
	result := app.Info(abci.RequestInfo{})
	assert.Equal("mockapp", result.GetData())
	// Should be 1 because we committed
	assert.Equal(int64(1), result.GetLastBlockHeight())
	hash1 := result.GetLastBlockAppHash()
	assert.NotNil(hash1)
	// Should == the first commit hash
	assert.Equal(c1.Data, hash1)

	// Call Query & Check the state of bobs account
	respQ := app.Query(abci.RequestQuery{Path: auth.AuthenticateQueryRoute, Data: bob.GetAddress()})
	assert.Equal(uint32(0), respQ.Code)

	ba := new(auth.BasicAccount)
	err := cdc.UnmarshalJSON(respQ.Value, ba)
	assert.Nil(err)
	assert.Equal(bob.GetAddress(), ba.Owner.Bytes())
	assert.Equal(uint32(10), ba.Balance)

	// Run validate
	// bob -> alice
	msg1 := auth.SendCoinMsg{Recipient: alice.GetAddress(), Amount: uint32(5)}
	tx, err := bob.SendMoney(cdc, msg1)
	assert.Nil(err)

	chtx := app.CheckTx(tx)
	assert.Equal(abci.ResponseCheckTx{Code: 0, Log: "Valid signature"}, chtx)

	bb := app.BeginBlock(abci.RequestBeginBlock{})
	assert.Equal(1, len(bb.Tags))
	assert.Equal([]byte("begin"), bb.Tags[0].Key)

	// Run Tx handler
	dtx := app.DeliverTx(tx)
	assert.Equal(abci.ResponseDeliverTx{Code: 0}, dtx)

	eb := app.EndBlock(abci.RequestEndBlock{})
	assert.Equal(1, len(eb.Tags))
	assert.Equal([]byte("end"), eb.Tags[0].Key)

	// Commit the new state
	commit := app.Commit()
	assert.NotNil(commit.Data)
	assert.NotEqual(c1.Data, commit.Data)

	// Check Alice's account & balance
	respQ = app.Query(abci.RequestQuery{Path: auth.AuthenticateQueryRoute, Data: alice.GetAddress()})
	assert.Equal(uint32(0), respQ.Code)
	al := new(auth.BasicAccount)
	err = cdc.UnmarshalJSON(respQ.Value, al)
	assert.Nil(err)
	assert.Equal(alice.GetAddress(), al.Owner.Bytes())
	assert.Equal(uint32(5), al.Balance)

	// Check bobs account again
	respQ = app.Query(abci.RequestQuery{Path: auth.AuthenticateQueryRoute, Data: bob.GetAddress()})
	assert.Equal(uint32(0), respQ.Code)

	ba = new(auth.BasicAccount)
	err = cdc.UnmarshalJSON(respQ.Value, ba)
	assert.Nil(err)
	assert.Equal(bob.GetAddress(), ba.Owner.Bytes())
	assert.Equal(uint32(5), ba.Balance)
}
