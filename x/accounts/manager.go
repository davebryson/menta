package accounts

// Manages storage access

import (
	"errors"

	"github.com/davebryson/menta/crypto"
	sdk "github.com/davebryson/menta/types"
	proto "github.com/golang/protobuf/proto"
)

var accountPrefix = []byte("menta/basicaccount")

// Prefixed account key
func getAccountKey(address []byte) []byte {
	return append(accountPrefix, address...)
}

func decodeAccount(raw []byte) (*Account, error) {
	acct := new(Account)
	err := proto.Unmarshal(raw, acct)
	if err != nil {
		return nil, err
	}
	return acct, nil
}

func encodeAccount(acct *Account) ([]byte, error) {
	return proto.Marshal(acct)
}

// LoadAccounts can used in initChain to bulk load genesis accounts from a json file
func LoadAccounts(store sdk.RWStore, accts []Account) error {
	for _, acct := range accts {
		err := SetAccount(store, acct)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAccount from state storage based on the account address
func GetAccount(store sdk.RWStore, address []byte) (*Account, error) {
	key := getAccountKey(address)
	acct, err := decodeAccount(store.Get(key))
	if err != nil {
		return nil, err
	}
	// Protobuf will return a struct with nil values if the []byte is empty!
	// versus an error
	if acct.GetPubkey() == nil {
		return nil, errors.New("Account not found")
	}

	return acct, nil
}

// SetAccount in the state store
func SetAccount(store sdk.RWStore, acct Account) error {
	pk, err := crypto.PublicKeyFromBytes(acct.Pubkey)
	if err != nil {
		return err
	}
	address := AddressFromPubKey(pk).Bytes()
	encodedAccount, err := encodeAccount(&acct)
	if err != nil {
		return err
	}
	store.Set(getAccountKey(address), encodedAccount)

	return nil
}
