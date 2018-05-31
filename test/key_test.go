package test

import (
	"encoding/hex"
	"fmt"
	"testing"

	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/go-crypto"
)

// Amino prefixes:
// Privatekey: a328891240
// Publickey: 1624DE6220
// Sig: 3da1db2a40
const (
	ADDRESS = "47b865865b12035ab5631d71e04ce08a4cb5bb63"
	SK      = "a328891240b8cdaa593df7789b293876d50ec0ff12e8e74805a349225291a5a1a9cec2ff248251708b5acc56fcf85d087306ffb31f80590d324fc1b1d718bac17b3d672763"
)

func TestKey(t *testing.T) {
	assert := assert.New(t)

	k1 := sdk.CreateKey()
	assert.NotNil(k1)
	skbits := k1.PrivateKey.Bytes()
	fmt.Printf("Before %x\n", k1)
	fmt.Printf("Adter %x\n", skbits)

	k2, err := sdk.KeyFromPrivateKeyBytes(skbits)
	assert.Nil(err)

	assert.Equal(k1.Address, k2.Address)

	// From hex
	userbits, err := hex.DecodeString(SK)
	assert.Nil(err)
	assert.NotNil(userbits)
	useraddress, err := hex.DecodeString(ADDRESS)
	assert.Nil(err)

	k3, err := sdk.KeyFromPrivateKeyBytes(userbits)
	assert.Nil(err)
	assert.NotNil(k3)
	assert.Equal(useraddress, k3.Address)

}

func TestSignature(t *testing.T) {
	assert := assert.New(t)
	k := sdk.KeyFromSecret([]byte("mysecret"))
	msg := crypto.Sha256([]byte("hello"))
	sig := k.PrivateKey.Sign(msg)
	r := k.PrivateKey.PubKey().VerifyBytes(msg, sig)
	assert.True(r)

	// All Data:
	//fmt.Printf("Private: %x\n", k.PrivateKey.Bytes())
	//fmt.Printf("Public: %x\n", k.PrivateKey.PubKey().Bytes())
	//fmt.Printf("Address: %x\n", k.Address)
	//fmt.Printf("Msg: %x\n", msg)
	//fmt.Printf("Sig for 'sha256(hello)': %x\n", sig.Bytes())

	// Test signature from menta-js keys
	sig2 := "3da1db2a403bcde685ccf1607ae9c786d3fe8ff087920905d34efa9a689dcc703a728e844465b23359e322d972a93d6d7c4b71c624734f3ee404bc0d437c9e5c5f98faa10b"
	msg2 := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	pub2 := "1624DE6220e62e8aa3acdede131789dbec920fcb6cdf801be69df416e0e3660f0e0a5ed82e"
	s, _ := hex.DecodeString(sig2)
	m, _ := hex.DecodeString(msg2)
	p, _ := hex.DecodeString(pub2)

	rt := sdk.Verify(m, s, p)
	assert.True(rt)
}
