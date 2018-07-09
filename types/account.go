package types

import (
	proto "github.com/golang/protobuf/proto"
	crypto "github.com/tendermint/tendermint/crypto"
)

type Account struct {
	Nonce   uint64
	Balance uint64
	PubKey  crypto.PubKey
}

// Create a new account based on a public key
func AccountFromPubKey(pubKey []byte) (*Account, error) {
	pk, err := crypto.PubKeyFromBytes(pubKey)
	if err != nil {
		return nil, err
	}
	return &Account{
		Nonce:   0,
		Balance: 0,
		PubKey:  pk,
	}, nil
}

func AccountFromBytes(raw []byte) (*Account, error) {
	var acctbits AccountBytes
	err := proto.Unmarshal(raw, &acctbits)
	if err != nil {
		return nil, err
	}

	pk, err := crypto.PubKeyFromBytes(acctbits.PubkeyBytes)
	if err != nil {
		return nil, err
	}

	return &Account{
		Nonce:   acctbits.Nonce,
		Balance: acctbits.Balance,
		PubKey:  pk,
	}, nil
}

func (acct *Account) Bytes() ([]byte, error) {
	return proto.Marshal(&AccountBytes{
		Nonce:       acct.Nonce,
		Balance:     acct.Balance,
		PubkeyBytes: acct.PubKey.Bytes(),
	})
}

func (acct *Account) Address() []byte {
	return acct.PubKey.Address().Bytes()
}
