package types

import (
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	OK uint32 = iota
	HandlerNotFound
	BadNonce
	BadSignature
	NoAccount
	BadTx
	NotFound
)

// ** Tag code from cosmos sdk **
type Tag = cmn.KVPair
type Tags cmn.KVPairs

func (t Tags) AppendTag(k string, v []byte) Tags {
	return append(t, MakeTag(k, v))
}

// Make a tag from a key and a value
func MakeTag(k string, v []byte) Tag {
	return Tag{Key: []byte(k), Value: v}
}

type Result struct {
	Code uint32 // Any non-zero code is an error
	Data []byte
	Log  string
	Tags Tags
}

func ResultError(code uint32, msg string) Result {
	return Result{
		Code: code,
		Log:  msg,
	}
}

func ErrorNoHandler() Result {
	return ResultError(HandlerNotFound, "Handler not found")
}

func ErrorBadNonce() Result {
	return ResultError(BadNonce, "Bad Nonce")
}

func ErrorBadSignature() Result {
	return ResultError(BadSignature, "Bad Signature")
}

func ErrorNoAccount() Result {
	return ResultError(NoAccount, "Account not found")
}

func ErrorBadTx() Result {
	return ResultError(BadTx, "Error decoding the transaction")
}
