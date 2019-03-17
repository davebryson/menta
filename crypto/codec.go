package crypto

import (
	"os"

	amino "github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}

// ShowAminoPrefix - From tendermint/crypto. Shows a nice table of the current crypto prefixes
func ShowAminoPrefix() {
	cdc.PrintTypes(os.Stdout)
}
