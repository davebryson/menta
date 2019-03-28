package authenticate

import (
	"fmt"

	sdk "github.com/davebryson/menta/types"
	"github.com/tendermint/tendermint/libs/common"
)

var _ sdk.Msg = (*SendCoinMsg)(nil)

// SendCoinMsg msg format for moving around coins
type SendCoinMsg struct {
	Recipient common.HexBytes
	Amount    uint32
}

// Route implements Msg
func (m SendCoinMsg) Route() string { return AuthenticateRoute }

// Type implements Msg
func (m SendCoinMsg) Type() string { return "" }

// ValidateBasic implements Msg
func (m SendCoinMsg) ValidateBasic() error {
	if m.Amount == 0 {
		return fmt.Errorf("Missing amount")
	}
	if m.Recipient == nil {
		return fmt.Errorf("Missing recipient")
	}
	return nil
}
