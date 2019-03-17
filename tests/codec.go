package test

import (
	"fmt"

	sdk "github.com/davebryson/menta/types"
	amino "github.com/tendermint/go-amino"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/libs/common"
)

var cdc = amino.NewCodec()
var _ sdk.Msg = (*HelloMsg)(nil)
var _ sdk.Msg = (*FunnyMoneyMsg)(nil)

// HelloMsg - msg test
type HelloMsg struct {
	Name string `json:"name"`
}

func (m HelloMsg) Route() string        { return "hello" }
func (m HelloMsg) Type() string         { return "info" }
func (m HelloMsg) ValidateBasic() error { return nil }

// FunnyMoneyMsg - app test
type FunnyMoneyMsg struct {
	Recipient common.HexBytes `json:"recipient"`
	Amount    uint32          `json:"amount"`
}

func (m FunnyMoneyMsg) Route() string        { return "funny" }
func (m FunnyMoneyMsg) Type() string         { return "exchange" }
func (m FunnyMoneyMsg) ValidateBasic() error { return nil }

// Storage models
type FunnyAcct struct {
	Owner   common.HexBytes
	Balance uint32
	PubKey  tmcrypto.PubKey
}

func (acct *FunnyAcct) Encode() ([]byte, error) {
	return cdc.MarshalBinaryBare(acct)
}

func DecodeAcct(accBits []byte) (*FunnyAcct, error) {
	senderAcc := new(FunnyAcct)
	err := cdc.UnmarshalBinaryBare(accBits, senderAcc)
	if err != nil {
		return nil, err
	}
	return senderAcc, nil
}

func (acct *FunnyAcct) Debit(amt uint32) error {
	if acct.Balance < amt {
		return fmt.Errorf("insufficient funds")
	}
	acct.Balance -= amt
	return nil
}
func (acct *FunnyAcct) Credit(amt uint32) {
	acct.Balance += amt
}

func init() {
	cryptoAmino.RegisterAmino(cdc)
	sdk.RegisterStandardTypes(cdc)
	cdc.RegisterConcrete(&HelloMsg{}, "our/hello", nil)
	cdc.RegisterConcrete(&FunnyMoneyMsg{}, "our/funnymoney", nil)
	cdc.RegisterConcrete(&FunnyAcct{}, "our/funnyacct", nil)
}
