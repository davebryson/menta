package types

import amino "github.com/tendermint/go-amino"

func RegisterStandardTypes(cdc *amino.Codec) {
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterInterface((*Tx)(nil), nil)
	cdc.RegisterConcrete(&StdTx{}, "menta/stdtx", nil)
	cdc.RegisterConcrete(&CommitInfo{}, "menta/state", nil)
}
