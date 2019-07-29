package types

import (
	"github.com/davebryson/menta/crypto"
	proto "github.com/golang/protobuf/proto"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

// DecodeTx returns a Tx from a []byte
func DecodeTx(raw []byte) (*Tx, error) {
	var tx Tx
	err := proto.Unmarshal(raw, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// EncodeTx returns a []byte or error
func EncodeTx(tx *Tx) ([]byte, error) {
	return proto.Marshal(tx)
}

// Hash the tx for signing
func (tx *Tx) hash() ([]byte, error) {
	bits, err := proto.Marshal(&Tx{
		Sender:  tx.Sender,
		Route:   tx.Route,
		Msgtype: tx.Msgtype,
		Msg:     tx.Msg,
		Nonce:   tx.Nonce,
	})
	if err != nil {
		return nil, err
	}
	hash := tmcrypto.Sha256(bits)
	return hash[:], nil
}

// Sign a transaction
func (tx *Tx) Sign(sk crypto.PrivateKeyEd25519) error {
	//tx.Sender = sk.PubKey().ToAddress().Bytes()
	msgHash, err := tx.hash()
	if err != nil {
		return err
	}
	tx.Sig = sk.Sign(msgHash)
	return nil
}

// Verify a Tx against a given public key
func (tx *Tx) Verify(pubKey crypto.PublicKeyEd25519) bool {
	msg, err := tx.hash()
	if err != nil {
		return false
	}
	return pubKey.Verify(msg, tx.Sig)
}
