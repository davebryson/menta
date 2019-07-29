package accounts

import (
	"testing"

	menta "github.com/davebryson/menta/app"
	"github.com/davebryson/menta/crypto"
	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestVerifyHandler(t *testing.T) {
	assert := assert.New(t)

	bobSk := crypto.GeneratePrivateKey()
	aliceSk := crypto.GeneratePrivateKey()

	app := menta.NewMockApp()
	app.OnInitChain(func(store sdk.RWStore, req abci.RequestInitChain) (resp abci.ResponseInitChain) {
		SetAccount(store, Account{Pubkey: bobSk.PubKey().Bytes()})
		return
	})
	app.OnValidateTx(VerifyAccount)
	app.OnTx("fakeroute", func(store sdk.RWStore, tx *sdk.Tx) sdk.Result {
		return sdk.Result{}
	})

	// Simulate
	app.InitChain(abci.RequestInitChain{})
	app.Commit()

	t1 := &sdk.Tx{Route: "fakeroute", Msg: []byte("hello")}
	t1.Sender = AddressFromPubKey(bobSk.PubKey()).Bytes()
	t1.Sign(bobSk)
	t1Bits, err := sdk.EncodeTx(t1)
	assert.Nil(err)

	// Good check
	resp := app.CheckTx(abci.RequestCheckTx{Tx: t1Bits})
	assert.Equal(uint32(0), resp.Code)

	// Bad route
	t1 = &sdk.Tx{Route: "badroute", Msg: []byte("hello")}
	t1.Sender = AddressFromPubKey(bobSk.PubKey()).Bytes()
	t1.Sign(bobSk)
	t1Bits, err = sdk.EncodeTx(t1)
	assert.Nil(err)

	resp = app.CheckTx(abci.RequestCheckTx{Tx: t1Bits})
	assert.True(resp.Code != 0)

	// Account not found
	t2 := &sdk.Tx{Route: "fakeroute", Msg: []byte("hello")}
	t2.Sender = AddressFromPubKey(aliceSk.PubKey()).Bytes()
	t2.Sign(aliceSk)
	assert.Equal(AddressFromPubKey(aliceSk.PubKey()).Bytes(), t2.Sender)
	t2Bits, err := sdk.EncodeTx(t2)
	assert.Nil(err)

	resp2 := app.CheckTx(abci.RequestCheckTx{Tx: t2Bits})
	assert.Equal(uint32(2), resp2.Code)

	// Bad signature
	t1 = &sdk.Tx{Route: "fakeroute", Msg: []byte("hello")}
	t1.Sender = AddressFromPubKey(bobSk.PubKey()).Bytes()
	t1.Sign(bobSk)
	t1.Sig[4] ^= byte(0x01) // Make signature bad
	t1Bits, err = sdk.EncodeTx(t1)
	assert.Nil(err)

	resp = app.CheckTx(abci.RequestCheckTx{Tx: t1Bits})
	assert.Equal(uint32(3), resp.Code)
}
