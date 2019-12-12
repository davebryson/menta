package types

import (
	"testing"

	"github.com/davebryson/menta/crypto"
	"github.com/stretchr/testify/assert"
)

func TestTx(t *testing.T) {
	assert := assert.New(t)

	bob := crypto.GeneratePrivateKey()

	tx := &SignedTransaction{Service: "one", Nonce: []byte("random"), Msg: []byte("hello")}
	tx.Sign(bob)

	txbits, err := EncodeTx(tx)
	assert.Nil(err)

	txBack, err := DecodeTx(txbits)
	assert.Nil(err)

	assert.Equal("one", txBack.Service)
	assert.Equal([]byte("random"), txBack.Nonce)
	assert.Equal([]byte("hello"), txBack.Msg)
	assert.True(txBack.Verify(bob.PubKey()))

	// Backwards from hex of private key
	bobSecretHex := bob.ToHex()
	bobBack, err := crypto.PrivateKeyFromHex(bobSecretHex)
	assert.Nil(err)
	assert.Equal(bob.PubKey().Bytes(), bobBack.PubKey().Bytes())

	// Backward from pubkey hex
	bobPubHex := bob.PubKey().ToHex()
	bob2, err := crypto.PublicKeyFromHex(bobPubHex)
	assert.Equal(bob.PubKey().Bytes(), bob2.Bytes())
}
