package types

import (
	"encoding/hex"
	"testing"

	"github.com/davebryson/menta/crypto"
	"github.com/stretchr/testify/assert"
)

// Test data generated from menta-js
var (
	jsMsgRoute        = "one"
	jsMsgData         = []byte("helloworld")
	jsWalletPrivKey   = "a7acc81ddb563e72d8cd2db7630ff63943afbe8b4db3be63f2460ede8395cbb1d2fca902543a126b65577add39f5826d443d1edfb558f9c8a2307fc8760fe378"
	jsWalletPubKey    = "d2fca902543a126b65577add39f5826d443d1edfb558f9c8a2307fc8760fe378"
	jsWalletAddress   = "4c8e740a9851aac6f9191f181449b1a1f2174a04"
	jsWalletSignedMsg = "0a036f6e651a0a68656c6c6f776f726c6422144c8e740a9851aac6f9191f181449b1a1f2174a042a0fe356f7f37736e1c71cdf8eb97fd7f73240d9ac0c6e395f359dc03e65d291ef1ab6ac3ef49e8cba70537c507e12beae4af053050d5a0145bd873dbb881fe2435281028392820c5449922581659b766db700"
)

func TestTx(t *testing.T) {
	assert := assert.New(t)

	bob := crypto.GeneratePrivateKey()

	tx := &Tx{Route: "one", Nonce: []byte("random"), Msg: []byte("hello")}
	tx.Sign(bob)

	txbits, err := EncodeTx(tx)
	assert.Nil(err)

	txBack, err := DecodeTx(txbits)
	assert.Nil(err)

	assert.Equal("one", txBack.Route)
	assert.Equal([]byte("random"), txBack.Nonce)
	assert.Equal([]byte("hello"), txBack.Msg)
	assert.True(txBack.Verify(bob.PubKey()))
	assert.Equal(bob.PubKey().ToAddress().Bytes(), txBack.Sender)

	// Backwards from hex of private key
	bobSecretHex := bob.ToHex()
	bobBack, err := crypto.PrivateKeyFromHex(bobSecretHex)
	assert.Nil(err)
	assert.Equal(bob.PubKey().Bytes(), bobBack.PubKey().Bytes())
	assert.Equal(bob.PubKey().ToAddress(), bobBack.PubKey().ToAddress())

	// Backward from pubkey hex
	bobPubHex := bob.PubKey().ToHex()
	bob2, err := crypto.PublicKeyFromHex(bobPubHex)
	assert.Equal(bob.PubKey().Bytes(), bob2.Bytes())
}

func TestWalletIntegration(t *testing.T) {
	assert := assert.New(t)

	sk, err := crypto.PrivateKeyFromHex(jsWalletPrivKey)
	assert.Nil(err)
	assert.NotNil(sk)

	bobPub := sk.PubKey().ToHex()
	assert.Equal(bobPub, jsWalletPubKey)
	bobAddy := sk.PubKey().ToAddress().ToHex()
	assert.Equal(bobAddy, jsWalletAddress)

	// Decode the msg into bytes
	bits, e := hex.DecodeString(jsWalletSignedMsg)
	assert.Nil(e)
	assert.NotNil(bits)

	// Decode Tx
	tx, e1 := DecodeTx(bits)
	assert.Nil(e1)

	assert.Equal(jsMsgRoute, tx.Route)
	assert.Equal(jsMsgData, tx.Msg)
	assert.True(tx.Verify(sk.PubKey()))
}
