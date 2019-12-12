package types

import (
	"github.com/davebryson/menta/crypto"
	proto "github.com/golang/protobuf/proto"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

// DecodeTx returns a Tx from a []byte
func DecodeTx(raw []byte) (*SignedTransaction, error) {
	var tx SignedTransaction
	err := proto.Unmarshal(raw, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// EncodeTx returns a []byte or error
func EncodeTx(tx *SignedTransaction) ([]byte, error) {
	return proto.Marshal(tx)
}

// Hash the tx for signing
func (tx *SignedTransaction) hashMsg() ([]byte, error) {
	bits, err := proto.Marshal(&SignedTransaction{
		Sender:  tx.Sender,
		Service: tx.Service,
		Msg:     tx.Msg,
		Msgid:   tx.Msgid,
		Nonce:   tx.Nonce,
	})
	if err != nil {
		return nil, err
	}
	hash := tmcrypto.Sha256(bits)
	return hash[:], nil
}

// Sign a transaction
func (tx *SignedTransaction) Sign(sk crypto.PrivateKeyEd25519) error {
	msgHash, err := tx.hashMsg()
	if err != nil {
		return err
	}
	tx.Sig = sk.Sign(msgHash)
	return nil
}

// Verify a Tx against a given public key
func (tx *SignedTransaction) Verify(pubKey crypto.PublicKeyEd25519) bool {
	msg, err := tx.hashMsg()
	if err != nil {
		return false
	}
	return pubKey.Verify(msg, tx.Sig)
}
