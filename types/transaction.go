package types

import (
	proto "github.com/golang/protobuf/proto"
	crypto "github.com/tendermint/go-crypto"
)

func TransactionFromBytes(raw []byte) (*Transaction, error) {
	var tx Transaction
	err := proto.Unmarshal(raw, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (tx *Transaction) Bytes() ([]byte, error) {
	return proto.Marshal(tx)
}

func (tx *Transaction) Hash() ([]byte, error) {
	bits, err := proto.Marshal(&Transaction{
		From:  tx.From,
		To:    tx.To,
		Nonce: tx.Nonce,
		Value: tx.Value,
		Call:  tx.Call,
		Data:  tx.Data,
	})

	if err != nil {
		return nil, err
	}
	return crypto.Sha256(bits), nil
}

func (tx *Transaction) Sign(key crypto.PrivKey) error {
	// Always set From to signer
	tx.From = key.PubKey().Address()

	msgHash, err := tx.Hash()
	if err != nil {
		return err
	}
	tx.Sig = key.Sign(msgHash).Bytes()
	return nil
}

func (tx *Transaction) Verify(pubKey crypto.PubKey) bool {
	hash, e0 := tx.Hash()
	if e0 != nil {
		return false
	}
	sig, e1 := crypto.SignatureFromBytes(tx.Sig)
	if e1 != nil {
		return false
	}
	return pubKey.VerifyBytes(hash, sig)
}
