package test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
)

// Amino prefixes:
// Privatekey: a328891040
// Publickey:  1624de6420
// Sig: 3da1db2a40
//  2031ea5340
const (
	ADDRESS = "2e3068298c64cc6095cd959d40e5549be7b2b04d"
	SK      = "a328891040f1db26ce33045d02c41a9888576705cd7fa3542afaa7fbcb2f5aceca43b07641d43302e2f82f7f5d42342c8f988ec5e17f1056feb179bf9803c35af390aa1a90"
)

func TestKey(t *testing.T) {
	assert := assert.New(t)

	k1 := sdk.CreateKey()
	assert.NotNil(k1)
	skbits := k1.PrivateKey.Bytes()

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

/*func TestSignature(t *testing.T) {
	assert := assert.New(t)
	k := sdk.KeyFromSecret([]byte("mysecret"))
	msg := crypto.Sha256([]byte("hello"))
	sig, err := k.PrivateKey.Sign(msg)
	assert.Nil(err)
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
}*/
