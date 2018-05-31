package types

import (
	"encoding/hex"
	"encoding/json"

	crypto "github.com/tendermint/go-crypto"
)

type Key struct {
	Address    []byte
	PrivateKey crypto.PrivKey
}

type keyAsJson struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privatekey"`
}

func Verify(msg, sig, pubkey []byte) bool {
	pk, err := crypto.PubKeyFromBytes(pubkey)
	if err != nil {
		return false
	}
	signature, err := crypto.SignatureFromBytes(sig)
	if err != nil {
		return false
	}
	return pk.VerifyBytes(msg, signature)
}

func KeyFromSecret(secret []byte) *Key {
	sk := crypto.GenPrivKeyEd25519FromSecret(secret)
	return &Key{
		Address:    sk.PubKey().Address(),
		PrivateKey: sk,
	}
}

func KeyFromPrivateKeyBytes(pk []byte) (*Key, error) {
	sk, err := crypto.PrivKeyFromBytes(pk)
	if err != nil {
		return nil, err
	}
	return &Key{
		Address:    sk.PubKey().Address(),
		PrivateKey: sk,
	}, nil
}

func CreateKey() *Key {
	sk := crypto.GenPrivKeyEd25519()
	return &Key{
		Address:    sk.PubKey().Address(),
		PrivateKey: sk,
	}
}

func (key *Key) ToJSON() ([]byte, error) {
	return json.Marshal(&keyAsJson{
		Address:    hex.EncodeToString(key.Address),
		PrivateKey: hex.EncodeToString(key.PrivateKey.Bytes()),
	})
}
