package codec

import (
	"os"

	amino "github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

type Codec = amino.Codec

func New() *Codec {
	return amino.NewCodec()
}

func RegisterCrypto(cdc *Codec) {
	cryptoAmino.RegisterAmino(cdc)
}

func ShowAminoPrefix() {
	cdc := New()
	RegisterCrypto(cdc)
	cdc.PrintTypes(os.Stdout)
}
