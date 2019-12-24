package counter

import (
	mcrypto "github.com/davebryson/menta/crypto"
	sdk "github.com/davebryson/menta/types"
)

type Wallet struct {
	secretKey mcrypto.PrivateKeyEd25519
}

func CreateWallet() Wallet {
	return Wallet{
		secretKey: mcrypto.GeneratePrivateKey(),
	}
}

func WalletFromSeed(seed string) Wallet {
	return Wallet{
		secretKey: mcrypto.PrivateKeyFromSecret([]byte(seed)),
	}

}

func (wallet Wallet) NewTx(val uint32) ([]byte, error) {
	encoded, err := NewCounter(val).Encode()
	if err != nil {
		return nil, err
	}
	t := &sdk.SignedTransaction{Service: ServiceName, Msg: encoded}
	t.Sign(wallet.secretKey)
	return sdk.EncodeTx(t)
}

func (wallet Wallet) PubKey() []byte {
	return wallet.secretKey.PubKey().Bytes()
}
