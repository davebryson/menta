package authenticate

import (
	"fmt"

	"github.com/davebryson/menta/codec"

	"github.com/davebryson/menta/store"
	sdk "github.com/davebryson/menta/types"
	"github.com/tendermint/tendermint/libs/common"
)

var accountPrefix = []byte("menta/basicaccount")

func getAccountKey(account []byte) []byte {
	return append(accountPrefix, account...)
}

type AccountManager struct {
	cdc *codec.Codec
}

func NewAccountManager(cdc *codec.Codec) AccountManager {
	return AccountManager{cdc: cdc}
}

func (manager AccountManager) AccountExists(ctx sdk.Context, address common.HexBytes) bool {
	acct := ctx.Db.Get(getAccountKey(address))
	return acct != nil
}

func (manager AccountManager) GetAccount(ctx sdk.Context, address common.HexBytes) (*BasicAccount, error) {
	if !manager.AccountExists(ctx, address) {
		return nil, fmt.Errorf("account doesn't exist")
	}
	acctbits := ctx.Db.Get(getAccountKey(address))
	account := new(BasicAccount)
	err := manager.cdc.UnmarshalBinaryBare(acctbits, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (manager AccountManager) QueryAccount(ctx store.QueryContext, address common.HexBytes) (*BasicAccount, error) {
	acctbits := ctx.Get(getAccountKey(address))
	account := new(BasicAccount)
	err := manager.cdc.UnmarshalBinaryBare(acctbits, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (manager AccountManager) SaveAccount(ctx sdk.Context, acct *BasicAccount) error {
	bits, err := manager.cdc.MarshalBinaryBare(acct)
	if err != nil {
		return err
	}
	ctx.Db.Set(getAccountKey(acct.Owner), bits)
	return nil
}
