package accounts

// Helper code to generate and load accounts in JSON

import (
	"encoding/json"

	"github.com/davebryson/menta/crypto"
)

// GenerateJSONAccounts will create a json list of private keys
func GenerateJSONAccounts(num int) ([]byte, error) {
	list := make([]string, 0)
	for i := 0; i < num; i++ {
		k := crypto.GeneratePrivateKey()
		skHex := k.ToHex()
		list = append(list, skHex)
	}
	data, err := json.MarshalIndent(list, "", " ")
	if err != nil {
		return nil, err
	}

	return data, nil
}

// LoadJSONAccounts loads a json list from above into []Account
func LoadJSONAccounts(data []byte) ([]Account, error) {
	list := make([]string, 0)
	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	accts := make([]Account, 0)
	for _, sk := range list {
		k, err := crypto.PrivateKeyFromHex(sk)
		if err != nil {
			return accts, err
		}
		accts = append(accts, Account{Pubkey: k.PubKey().Bytes(), Nonce: 0})
	}

	return accts, nil
}
