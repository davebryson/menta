package accounts

import (
	"encoding/hex"
	"errors"

	mentakeys "github.com/davebryson/menta/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const addressSize = 20

type Address [addressSize]byte

func AddressFromPubKey(pubKey mentakeys.PublicKeyEd25519) Address {
	results := tmhash.SumTruncated(pubKey[:])
	var addressBits [addressSize]byte
	copy(addressBits[:], results)
	return Address(addressBits)
}

func AddressFromHex(h string) (Address, error) {
	bits, err := hex.DecodeString(h)
	if err != nil {
		return Address{}, err
	}
	if len(bits) != addressSize {
		return Address{}, errors.New("Address: not a valid address size")
	}
	var aBits [addressSize]byte
	copy(aBits[:], bits)
	return Address(aBits), nil
}

// Bytes returns the address as bytes
func (addy Address) Bytes() []byte {
	return addy[:]
}

// ToHex returns the address as a hex
func (addy Address) ToHex() string {
	return hex.EncodeToString(addy[:])
}
