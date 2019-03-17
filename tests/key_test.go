package test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestCryptoBasics(t *testing.T) {
	assert := assert.New(t)

	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()

	fmt.Printf("Addy: %v\n", pubKey.Address())
	//
	fmt.Printf("PubKey: %v\n", pubKey.Bytes())
	r := hex.EncodeToString(pubKey.Bytes())
	fmt.Printf("PubKey: %v\n", r)

	raw := pubKey.(ed25519.PubKeyEd25519)
	fmt.Printf("W/OPrefix PubKey: %v\n", raw.String())

	msg := crypto.CRandBytes(128)
	sig, err := privKey.Sign(msg)
	assert.Nil(err)

	// Test the signature
	assert.True(pubKey.VerifyBytes(msg, sig))

	// Mutate the signature, just one bit.
	// TODO: Replace this with a much better fuzzer, tendermint/ed25519/issues/10
	sig[7] ^= byte(0x01)

	assert.False(pubKey.VerifyBytes(msg, sig))
}
