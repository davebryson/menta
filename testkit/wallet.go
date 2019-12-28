package testkit

import (
	mcrypto "github.com/davebryson/menta/crypto"
	sdk "github.com/davebryson/menta/types"
	"github.com/gogo/protobuf/proto"
)

// Wallet provides a way to generate and sign transactions
type Wallet struct {
	secretKey mcrypto.PrivateKeyEd25519
}

// RandomWallet creates a new Wallet
func RandomWallet() Wallet {
	return Wallet{
		secretKey: mcrypto.GeneratePrivateKey(),
	}
}

// WalletFromSeed create a wallet based on the seed. This is a good way to create
// a predictable wallet with the same private key
func WalletFromSeed(seed string) Wallet {
	return Wallet{
		secretKey: mcrypto.PrivateKeyFromSecret([]byte(seed)),
	}

}

// CreateTx generates and signs a transaction return it as encoded bytes
func (wallet Wallet) CreateTx(serviceName string, msgid uint32, message proto.Message) ([]byte, error) {
	encoded, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	t := &sdk.SignedTransaction{Service: serviceName, Msgid: msgid, Msg: encoded}
	t.Sign(wallet.secretKey)
	return sdk.EncodeTx(t)
}

// PubKey returns the publickey for the wallet as bytes
func (wallet Wallet) PubKey() []byte {
	return wallet.secretKey.PubKey().Bytes()
}
