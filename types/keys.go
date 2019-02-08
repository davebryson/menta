package types

import (
	"crypto/sha256"
	"errors"

	signer "github.com/kevinburke/nacl/sign"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

type KeyPair struct {
	Public  signer.PublicKey
	Private signer.PrivateKey
	Address []byte
}

func GenerateKeyPair() (*KeyPair, error) {
	pub, priv, err := signer.Keypair(tmcrypto.CReader())
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		Public:  pub,
		Private: priv,
		Address: GetKeyAddress(pub),
	}, nil
}

func KeyPairFromPrivateKey(bits []byte) (*KeyPair, error) {
	if len(bits) != 64 {
		return nil, errors.New("Key must be 64 bits")
	}
	priv := signer.PrivateKey(bits)
	pub := priv.Public().(signer.PublicKey)
	return &KeyPair{
		Public:  pub,
		Private: priv,
		Address: GetKeyAddress(pub),
	}, nil
}

// GetKeyAddress returns the first 20 bytes of the sha256 of the public key
func GetKeyAddress(pub signer.PublicKey) []byte {
	hash := sha256.Sum256(pub[:])
	return hash[:20]
}
