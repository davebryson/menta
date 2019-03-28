package authenticate

import (
	"github.com/davebryson/menta/codec"
	sdk "github.com/davebryson/menta/types"
)

// RegisterAuthenticateTypes for codec
func RegisterAuthenticateTypes(cdc *codec.Codec) {
	sdk.RegisterStandardTypes(cdc)
	cdc.RegisterConcrete(&SendCoinMsg{}, "menta/auth/sendcoinmsg", nil)
}
