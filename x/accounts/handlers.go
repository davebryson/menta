package accounts

// Menta app handlers

import (
	"bytes"

	"github.com/davebryson/menta/crypto"
	sdk "github.com/davebryson/menta/types"
)

var (
	AccountNotFound  = sdk.ResultError(2, "Account not found")
	BadSignature     = sdk.ResultError(3, "Bad signature")
	FailedToSave     = sdk.ResultError(5, "Failed to save account")
	BadAccountPubKey = sdk.ResultError(6, "Bad account public key")
	PubKeyNoMatch    = sdk.ResultError(7, "Publickey address does not match the sender address")
)

// VerifyAccount handler is commonly used in checkTx to verify the sender
// Note: this is a very basic implementation.  It doesn't properly manage a nonce
func VerifyAccount(ctx sdk.Context) sdk.Result {
	tx := ctx.Tx

	// 1. Get the senders account
	acct, err := GetAccount(ctx, tx.Sender)
	if err != nil || acct == nil {
		return AccountNotFound
	}

	// 2. Convert acct pubkey bytes to PublicKey
	pubKey, err := crypto.PublicKeyFromBytes(acct.Pubkey)
	if err != nil {
		return BadAccountPubKey
	}

	// 3. Verify the account pubkey address matches the sender
	if !bytes.Equal(pubKey.ToAddress().Bytes(), tx.Sender) {
		return PubKeyNoMatch
	}

	// 4. Verify the signature
	goodSig := tx.Verify(pubKey)
	if !goodSig {
		return BadSignature
	}

	// Good to go...
	return sdk.Result{}
}
