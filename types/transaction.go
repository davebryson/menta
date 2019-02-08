package types

// Wrapper Tx
import (
	"crypto/sha256"

	proto "github.com/golang/protobuf/proto"
	signer "github.com/kevinburke/nacl/sign"
)

// TransactionFromBytes decodes a []byte to a Transaction
func TransactionFromBytes(raw []byte) (*Transaction, error) {
	var tx Transaction
	err := proto.Unmarshal(raw, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// ToBytes returns a serialized Tx
func (tx *Transaction) ToBytes() ([]byte, error) {
	return proto.Marshal(tx)
}

// Hash returns a sha256 hash of the Tx contents
func (tx *Transaction) Hash() ([]byte, error) {
	bits, err := proto.Marshal(&Transaction{
		Route:  tx.Route,
		Action: tx.Action,
		Nonce:  tx.Nonce,
		Data:   tx.Data,
	})

	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(bits)
	return hash[:], nil
}

// Sign a Transaction given a KeyPair
func (tx *Transaction) Sign(key *KeyPair) error {
	// Always set From to signer
	tx.Sender = key.Address

	msgHash, err := tx.Hash()
	if err != nil {
		return err
	}

	tx.Sig = signer.Sign(msgHash, key.Private)

	return nil
}

// Verify the transaction is valid given a public key
func (tx *Transaction) Verify(pubKey signer.PublicKey) bool {
	return signer.Verify(tx.Sig, pubKey)
}
