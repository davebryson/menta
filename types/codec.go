package types

import "github.com/davebryson/menta/codec"

func RegisterStandardTypes(cdc *codec.Codec) {
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterInterface((*Tx)(nil), nil)
	cdc.RegisterConcrete(&StdTx{}, "menta/stdtx", nil)
	cdc.RegisterConcrete(&CommitInfo{}, "menta/state", nil)
}
