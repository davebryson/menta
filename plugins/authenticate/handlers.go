package authenticate

import (
	"encoding/json"

	"github.com/tendermint/tendermint/libs/common"

	"github.com/davebryson/menta/store"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	AuthenticateRoute      = "/authenticate"
	AuthenticateQueryRoute = "/authenticate/query"
)

var (
	MissingStdTxResult      = sdk.ResultError(1, "Expected a stdTx")
	AccountNotFoundResult   = sdk.ResultError(2, "Account not found")
	BadSignatureResult      = sdk.ResultError(3, "Bad signature")
	InsufficientFundsResult = sdk.ResultError(4, "Insufficient funds")
	FailedToSaveResult      = sdk.ResultError(5, "Failed to save account")
)

// VerifyAccountHandler checks a tx for the proper signature.
// Intended for use in checkTx. Any Tx sent to the app will
// run through this handler
func VerifyAccountHandler(am AccountManager) sdk.TxHandler {
	return func(ctx sdk.Context) sdk.Result {
		stdTx, ok := ctx.Tx.(sdk.StdTx)
		if !ok {
			return MissingStdTxResult
		}

		// 1. Get the sender's account
		// 2. Make the msg for signature verification
		// 3. Verify the signature
		senderAcct, err := am.GetAccount(ctx, stdTx.GetSigner())
		if err != nil {
			return AccountNotFoundResult
		}

		// Get a sorted json byte representation of the msg
		sigmsg, err := json.Marshal(stdTx.GetMsg())

		// Get the actual signature
		signature := stdTx.GetSignBytes()
		if !senderAcct.PubKey.VerifyBytes(sigmsg, signature) {
			return BadSignatureResult
		}
		return sdk.Result{Log: "Valid signature"}
	}
}

// AccountQueryHandler is a simple query handler.  Will retursn an account in JSON to the caller
func AccountQueryHandler(am AccountManager) store.QueryHandler {
	return func(key []byte, ctx store.QueryContext) abci.ResponseQuery {
		res := abci.ResponseQuery{}
		res.Code = 10 // fail

		acct, err := am.QueryAccount(ctx, key)
		if err != nil {
			res.Log = err.Error()
			return res
		}
		jsonbits, err := am.cdc.MarshalJSON(acct)
		if err != nil {
			res.Log = err.Error()
			return res
		}

		res.Code = 0
		res.Value = jsonbits
		return res
	}
}

// AccountTxHandler can help exchange 'coin' among users.
func AccountTxHandler(am AccountManager) sdk.TxHandler {
	return func(ctx sdk.Context) sdk.Result {
		switch m := ctx.Tx.GetMsg().(type) {
		case *SendCoinMsg:
			return handleCoinExchange(ctx, am, ctx.Tx.GetSigner(), m)
		default:
			return sdk.ResultError(10, "Unrecognized msg")
		}
	}
}

func handleCoinExchange(ctx sdk.Context, am AccountManager, sender common.HexBytes, msg *SendCoinMsg) sdk.Result {
	senderAcct, err := am.GetAccount(ctx, sender)
	if err != nil {
		return AccountNotFoundResult
	}
	recipAccount, err := am.GetAccount(ctx, msg.Recipient)
	if err != nil {
		return AccountNotFoundResult
	}

	if senderAcct.Balance < msg.Amount {
		return InsufficientFundsResult
	}

	senderAcct.Debit(msg.Amount)
	recipAccount.Credit(msg.Amount)

	err = am.SaveAccount(ctx, senderAcct)
	if err != nil {
		return FailedToSaveResult
	}
	err = am.SaveAccount(ctx, recipAccount)
	if err != nil {
		return FailedToSaveResult
	}

	return sdk.Result{}
}
