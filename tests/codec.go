package test

import (
	"github.com/davebryson/menta/codec"
	auth "github.com/davebryson/menta/plugins/authenticate"
	sdk "github.com/davebryson/menta/types"
)

var cdc = codec.New()

var _ sdk.Msg = (*HelloMsg)(nil)

// HelloMsg - msg test
type HelloMsg struct {
	Name string `json:"name"`
}

func (m HelloMsg) Route() string        { return "hello" }
func (m HelloMsg) Type() string         { return "info" }
func (m HelloMsg) ValidateBasic() error { return nil }

func init() {
	codec.RegisterCrypto(cdc)
	auth.RegisterAuthenticateTypes(cdc)
	cdc.RegisterConcrete(&HelloMsg{}, "our/hello", nil)
}
