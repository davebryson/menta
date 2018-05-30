package test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
)

const (
	ADDRESS = "47b865865b12035ab5631d71e04ce08a4cb5bb63"
	SK      = "a328891240b8cdaa593df7789b293876d50ec0ff12e8e74805a349225291a5a1a9cec2ff248251708b5acc56fcf85d087306ffb31f80590d324fc1b1d718bac17b3d672763"
)

func TestKey(t *testing.T) {
	assert := assert.New(t)

	k1 := sdk.CreateKey()
	assert.NotNil(k1)
	skbits := k1.PrivateKey.Bytes()

	k2, err := sdk.KeyFromPrivateKeyBytes(skbits)
	assert.Nil(err)

	assert.Equal(k1.Address, k2.Address)

	// From hex
	userbits, err := hex.DecodeString(SK)
	assert.Nil(err)
	assert.NotNil(userbits)
	useraddress, err := hex.DecodeString(ADDRESS)
	assert.Nil(err)

	k3, err := sdk.KeyFromPrivateKeyBytes(userbits)
	assert.Nil(err)
	assert.NotNil(k3)
	assert.Equal(useraddress, k3.Address)

}
