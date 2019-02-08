package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	assert := assert.New(t)

	keys, err := GenerateKeyPair()
	assert.Nil(err)

	tx := &Transaction{}
	tx.Nonce = uint64(10)
	tx.Data = []byte("thepayload")

	err = tx.Sign(keys)
	assert.Nil(err)
	bits, err := tx.ToBytes()
	assert.Nil(err)

	tx1, e := TransactionFromBytes(bits)
	assert.Nil(e)

	assert.True(tx1.Verify(keys.Public))
}
