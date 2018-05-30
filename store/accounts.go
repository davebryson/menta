package store

import (
	"fmt"
)

var accountKeyPrefix = "acct/%x"

func cacheAccountKey(address []byte) string {
	return fmt.Sprintf(accountKeyPrefix, address)
}
