package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

func TestBasics(t *testing.T) {
	assert := assert.New(t)

	sk := GeneratePrivateKey()
	pk := sk.PubKey()

	// Sign and Verify
	msg := tmcrypto.Sha256([]byte("hello there"))
	sig := sk.Sign(msg)
	assert.Equal(64, len(sig))
	assert.True(pk.Verify(msg, sig))

	// Making bad...
	sig[7] ^= byte(0x01)
	assert.Equal(64, len(sig))
	assert.False(pk.Verify(msg, sig))

	// From hex
	skhex := sk.ToHex()
	sk1, err := PrivateKeyFromHex(skhex)
	assert.Nil(err)
	sig1 := sk1.Sign(msg)
	assert.Equal(64, len(sig1))
	assert.True(pk.Verify(msg, sig1))
}
