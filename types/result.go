package types

// Helper for returning results from check/deliver calls
import (
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	// OK - all is bueno with the executed Tx. Any non-zero code is an error
	OK uint32 = iota
	// HandlerNotFound - yep...we couldn't find it
	HandlerNotFound
	// BadTx - no bueno, couldn't decode it or something like that
	BadTx
	// NotFound - catch all
	NotFound
	// BadQuery - in store query
	BadQuery
)

// Result is it returned from a menta app TxHandler
// By default 'Code' will be zero which mean 'Ok' to tendermint
type Result struct {
	Code uint32 // Any non-zero code is an error
	Data []byte
	Log  string
}

// Tag defined in tendermint ... common/types.proto
type Tag = cmn.KVPair

// Tags defined in tendermint ... common/kvpairs.go = []cmn.KVPairs
type Tags cmn.KVPairs

// ResultError is returned on an error with a non-zero code
func ResultError(code uint32, log string) Result {
	return Result{
		Code: code,
		Log:  log,
	}
}

// ErrorNoHandler is returned when menta can't find a handler for a given route
func ErrorNoHandler() Result {
	return ResultError(HandlerNotFound, "Handler not found")
}

// ErrorBadTx is returned when menta can't deserialize a Tx
func ErrorBadTx() Result {
	return ResultError(BadTx, "Error decoding the transaction")
}
