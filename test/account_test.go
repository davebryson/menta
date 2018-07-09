package test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/tendermint/crypto"
)

const (
	ADDY   = "afe25dec1002986616c5f00a942e3b623e2e4f42"
	PUB    = "1624de64207559c028a191cc9e7edc41dabb8363c8ba0027ebbdffcd1749246b7d7800131f"
	SECRET = "helloworld"
	MSG    = "heythere"
)

// NOTE: amino pubkey prefix: 1624de6420

func TestAccountStuff(t *testing.T) {
	assert := assert.New(t)

	s, e0 := hex.DecodeString(PUB)
	assert.Nil(e0)

	a, e := types.AccountFromPubKey([]byte(s))
	assert.Nil(e)
	assert.NotNil(a)

	a.Nonce = uint64(1)
	a.Balance = uint64(10)
	assert.Equal(ADDY, hex.EncodeToString(a.Address()))

	raw, err := a.Bytes()
	assert.Nil(err)
	back, err := types.AccountFromBytes(raw)
	assert.Nil(err)

	assert.Nil(err)
	assert.Equal(a.Address(), back.Address())
	assert.Equal(a.PubKey, back.PubKey)
	assert.Equal(a.Nonce, back.Nonce)
	assert.Equal(a.Balance, back.Balance)
}

func TestAcctAndSigs(t *testing.T) {
	assert := assert.New(t)

	sk := crypto.GenPrivKeyEd25519FromSecret([]byte(SECRET))
	account, err := types.AccountFromPubKey(sk.PubKey().Bytes())
	assert.Nil(err)
	assert.NotNil(account)

	raw, err := account.Bytes()
	assert.Nil(err)
	assert.NotNil(raw)
	back, err := types.AccountFromBytes(raw)

	// Sign something ...
	//sigBytes := sk.Sign([]byte(MSG)).Bytes()
	sig, e := sk.Sign([]byte(MSG))
	assert.Nil(e)
	signature, err := crypto.SignatureFromBytes(sig.Bytes())

	fmt.Printf("Sig: %s", hex.EncodeToString(sig.Bytes()))
	assert.Nil(err)
	assert.NotNil(signature)

	// Verify with account
	result := back.PubKey.VerifyBytes([]byte(MSG), signature)
	assert.True(result)
}
