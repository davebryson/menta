package crypto

import (
	"encoding/hex"
	"errors"
	"io"

	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"golang.org/x/crypto/ed25519"
)

// NOTE: Some of this code is adapted from the Tendermint crypto library

const (
	// PublicKeySize is the size, in bytes, of public keys as used in this package.
	PublicKeySize = 32
	// PrivateKeySize is the size, in bytes, of private keys as used in this package.
	PrivateKeySize = 64
	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize = 64
	// AddressSize is the size in bytes of an address extracted from the public key
	AddressSize = 20
)

// PrivateKeyEd25519 is the private key container
type PrivateKeyEd25519 [PrivateKeySize]byte

// PublicKeyEd25519 is the public key container
type PublicKeyEd25519 [PublicKeySize]byte

// Signature is the signature container
type Signature [SignatureSize]byte

// Address is the address container
type Address [AddressSize]byte

// GeneratePrivateKey generates a new private key
func GeneratePrivateKey() PrivateKeyEd25519 {
	seed := make([]byte, 32)
	_, err := io.ReadFull(tmcrypto.CReader(), seed[:])
	if err != nil {
		panic(err)
	}
	privKey := ed25519.NewKeyFromSeed(seed)
	var privKeyEd PrivateKeyEd25519
	copy(privKeyEd[:], privKey)
	return privKeyEd
}

// PrivateKeyFromSecret generates a private key from a given secret
func PrivateKeyFromSecret(secret []byte) PrivateKeyEd25519 {
	seed := tmcrypto.Sha256(secret)
	privKey := ed25519.NewKeyFromSeed(seed)
	var privKeyEd PrivateKeyEd25519
	copy(privKeyEd[:], privKey)
	return privKeyEd
}

// PrivateKeyFromHex decodes a hexified private key into PrivateKeyEd25519
func PrivateKeyFromHex(h string) (PrivateKeyEd25519, error) {
	bits, err := hex.DecodeString(h)
	if err != nil {
		return PrivateKeyEd25519{}, err
	}
	if len(bits) != PrivateKeySize {
		return PrivateKeyEd25519{}, errors.New("privatekey: not a valid private key size")
	}

	var pkBits [PrivateKeySize]byte
	copy(pkBits[:], bits)
	return PrivateKeyEd25519(pkBits), nil
}

// Sign a message
func (privKey PrivateKeyEd25519) Sign(msg []byte) []byte {
	signatureBytes := ed25519.Sign(privKey[:], msg)
	return signatureBytes[:]
}

// PubKey returns the public key for this private key
func (privKey PrivateKeyEd25519) PubKey() PublicKeyEd25519 {
	privKeyBytes := [PrivateKeySize]byte(privKey)
	initialized := false
	// If the latter 32 bytes of the privkey are all zero, compute the pubkey
	// otherwise privkey is initialized and we can use the cached value inside
	// of the private key.
	for _, v := range privKeyBytes[32:] {
		if v != 0 {
			initialized = true
			break
		}
	}
	if !initialized {
		panic("Expected PrivKeyEd25519 to include concatenated pubkey bytes")
	}
	var pubkeyBytes [PublicKeySize]byte
	copy(pubkeyBytes[:], privKeyBytes[32:])
	return PublicKeyEd25519(pubkeyBytes)
}

// Bytes return the the private key as bytes
func (privKey PrivateKeyEd25519) Bytes() []byte {
	return privKey[:]
}

// ToHex returns the private key as a hex value
func (privKey PrivateKeyEd25519) ToHex() string {
	return hex.EncodeToString(privKey[:])
}

// ---- PublicKey ----

// PublicKeyFromHex decodes a hex version of the public key into PublicKeyEd25519
func PublicKeyFromHex(h string) (PublicKeyEd25519, error) {
	bits, err := hex.DecodeString(h)
	if err != nil {
		return PublicKeyEd25519{}, err
	}
	if len(bits) != PublicKeySize {
		return PublicKeyEd25519{}, errors.New("publickey: not a valid public key size")
	}

	var pkBits [PublicKeySize]byte
	copy(pkBits[:], bits)
	return PublicKeyEd25519(pkBits), nil
}

// ToAddress returns the associated address for the the key
func (pubKey PublicKeyEd25519) ToAddress() Address {
	results := tmhash.SumTruncated(pubKey[:])
	var addressBits [AddressSize]byte
	copy(addressBits[:], results)
	return Address(addressBits)
}

// Verify a signature and message
func (pubKey PublicKeyEd25519) Verify(msg []byte, sig []byte) bool {
	// make sure we use the same algorithm to sign
	if len(sig) != SignatureSize {
		return false
	}
	return ed25519.Verify(pubKey[:], msg, sig)
}

// Bytes returns the public key as bytes
func (pubKey PublicKeyEd25519) Bytes() []byte {
	return pubKey[:]
}

// ToHex returns the public key as a hex
func (pubKey PublicKeyEd25519) ToHex() string {
	return hex.EncodeToString(pubKey[:])
}

// ---- Address ----

// AddressFromHex decodes a hex version of the address into an Address
func AddressFromHex(h string) (Address, error) {
	bits, err := hex.DecodeString(h)
	if err != nil {
		return Address{}, err
	}
	if len(bits) != AddressSize {
		return Address{}, errors.New("address: not a valid address size")
	}
	var aBits [AddressSize]byte
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
