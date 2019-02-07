package test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
)

func TestBasics(t *testing.T) {
	assert := assert.New(t)

	tx := &sdk.Transaction{}
	tx.Nonce = uint64(10)

	bits, err := tx.Bytes()
	assert.Nil(err)
	assert.NotNil(bits)

	r, e := sdk.TransactionFromBytes(bits)
	assert.Nil(e)
	assert.NotNil(r)
	assert.Equal(uint64(10), r.Nonce)
}

func TestFullTx(t *testing.T) {
	assert := assert.New(t)

	tx := &sdk.Transaction{}
	tx.Nonce = uint64(10)
	tx.Value = uint64(2)
	tx.Call = "hello"
	tx.Data = []byte("thepayload")
	//err := tx.Sign(k.PrivateKey)
	//assert.Nil(err)

	// Serialize
	txbits, err := tx.Bytes()
	assert.Nil(err)
	assert.NotNil(txbits)

	// Thaw
	tx2, err := sdk.TransactionFromBytes(txbits)
	assert.Nil(err)

	// Verify
	//ok := tx2.Verify(k.PrivateKey.PubKey())
	//assert.True(ok)

	// Check contents
	assert.Equal(uint64(10), tx2.Nonce)
	assert.Equal(uint64(2), tx2.Value)
	assert.Equal("hello", tx2.Call)
	assert.Equal([]byte("thepayload"), tx2.Data)
}

func TestDecodeFromApi(t *testing.T) {
	assert := assert.New(t)

	// encoded tx from menta-js client
	raw, err := hex.DecodeString("18012a0548656c6c6f")
	assert.Nil(err)

	tx, err := sdk.TransactionFromBytes(raw)
	assert.Nil(err)
	assert.Equal("Hello", tx.Call)
	assert.Equal(uint64(1), tx.Nonce)

}
