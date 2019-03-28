package authenticate

import (
	"fmt"

	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/common"
)

// BasicAccount stores account information
type BasicAccount struct {
	Owner   common.HexBytes
	Balance uint32
	PubKey  tmcrypto.PubKey
}

// NewAccount returns a instance of a BasicAccount
func NewAccount(pubKey tmcrypto.PubKey) *BasicAccount {
	return &BasicAccount{
		Owner:   pubKey.Address(),
		Balance: uint32(0),
		PubKey:  pubKey,
	}
}

// Debit the account by the given amount
func (acct *BasicAccount) Debit(amt uint32) error {
	if acct.Balance < amt {
		return fmt.Errorf("insufficient funds")
	}
	acct.Balance -= amt
	return nil
}

// Credit the account for the given ammount
func (acct *BasicAccount) Credit(amt uint32) {
	acct.Balance += amt
}
