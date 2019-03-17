package test

import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/go-amino"

	mentapp "github.com/davebryson/menta/app"
	"github.com/davebryson/menta/store"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

const FunnyMoneyRouteKey = "funny"

// Simple wallet for sending txs
type FunnyMoneyWallet struct {
	secretKey tmcrypto.PrivKey
}

func NewWallet() FunnyMoneyWallet {
	return FunnyMoneyWallet{
		secretKey: ed25519.GenPrivKey(),
	}
}

func (wallet FunnyMoneyWallet) GetAddress() []byte {
	return wallet.secretKey.PubKey().Address()
}

func (wallet FunnyMoneyWallet) GetPubKey() tmcrypto.PubKey {
	return wallet.secretKey.PubKey()
}

func (wallet FunnyMoneyWallet) SendMoney(msg FunnyMoneyMsg) ([]byte, error) {
	// json encoding automatically sorts the keys
	sigmsg, err := json.Marshal(msg)
	sig, err := wallet.secretKey.Sign(sigmsg)
	if err != nil {
		return nil, err
	}
	tx := sdk.StdTx{}
	tx.Signer = wallet.GetAddress()
	tx.Signature = sig
	tx.Msg = msg
	return cdc.MarshalJSON(tx)
}

// Future auth handler
func MoneyExchangeCheckTx(ctx sdk.Context) sdk.Result {
	stdTx, ok := ctx.Tx.(sdk.StdTx)
	if !ok {
		return sdk.ResultError(1, "expected a std tx")
	}

	// 1. Get the sender's account
	// 2. Make the msg for signature verification
	// 3. Verify the signature
	senderAcctBits := ctx.Db.Get(stdTx.GetSigner())
	if senderAcctBits == nil {
		return sdk.ResultError(1, "account not found")
	}

	senderAcc, err := DecodeAcct(senderAcctBits)
	// First thaw and check the sender's account
	if err != nil {
		return sdk.ResultError(1, "problem unfreezing account")
	}

	// Get a sorted json byte representation of the msg
	sigmsg, err := json.Marshal(stdTx.GetMsg())

	// Get the actual signature
	signature := stdTx.GetSignBytes()
	// Get the senders account and decode it

	if !senderAcc.PubKey.VerifyBytes(sigmsg, signature) {
		return sdk.ResultError(10, "Bad signature on message")
	}
	return sdk.Result{}
}

// Handler
func MoneyExchange(ctx sdk.Context) sdk.Result {
	switch m := ctx.Tx.GetMsg().(type) {
	case *FunnyMoneyMsg:
		// get the sender's account (serialized)
		senderAcctBits := ctx.Db.Get(ctx.Tx.GetSigner())
		if senderAcctBits == nil {
			return sdk.ResultError(1, "account not found")
		}

		// First thaw and check the sender's account
		senderAcc, err := DecodeAcct(senderAcctBits)
		// First thaw and check the sender's account
		if err != nil {
			return sdk.ResultError(1, "problem unfreezing account")
		}

		if senderAcc.Balance < m.Amount {
			return sdk.ResultError(1, "insufficient fund!")
		}

		// Now get or create the recipients account
		recAcctBits := ctx.Db.Get(m.Recipient)
		recAcc := new(FunnyAcct)
		if recAcctBits == nil {
			recAcc.Owner = m.Recipient
			recAcc.Balance = 0
		} else {
			recAcc, err = DecodeAcct(recAcctBits)
			if err != nil {
				return sdk.ResultError(1, "problem unfreezing account")
			}
		}

		// Debit/Credit accounts
		senderAcc.Debit(m.Amount)
		recAcc.Credit(m.Amount)

		// Serialize both accounts to storage
		sbits, err := senderAcc.Encode()
		//sbits, err := cdc.MarshalBinaryBare(senderAcc)
		if err != nil {
			fmt.Printf("ERROR: %v", err)
			return sdk.ResultError(1, "sender problem saving...")
		}
		ctx.Db.Set(senderAcc.Owner, sbits)

		rbits, err := recAcc.Encode()
		if err != nil {
			return sdk.ResultError(1, "recp problem saving...")
		}
		ctx.Db.Set(recAcc.Owner, rbits)
		return sdk.Result{
			Code: 0, Log: "xfer",
		}
	default:
		return sdk.ResultError(2, "unknown msg")
	}

}

// Query Handler
func QueryAccount(key []byte, ctx store.QueryContext) abci.ResponseQuery {
	acctAddress := key // Only to be explicit here - the key is an acct address
	res := abci.ResponseQuery{}

	res.Code = 10 // fail
	raw := ctx.Get(acctAddress)
	if raw == nil {
		res.Log = "Account not found"
	}

	acct, err := DecodeAcct(raw)
	if err != nil {
		res.Log = "Problem decoding..."
	}

	jsonbits, err := cdc.MarshalJSON(acct)
	if err != nil {
		res.Log = "Problem making json..."
	}

	res.Code = 0
	res.Value = jsonbits
	return res
}

// App
func createApp(bob, alice FunnyMoneyWallet) *mentapp.MentaApp {
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	sdk.RegisterStandardTypes(cdc)
	cdc.RegisterConcrete(&HelloMsg{}, "our/hello", nil)
	cdc.RegisterConcrete(&FunnyMoneyMsg{}, "our/funnymoney", nil)
	cdc.RegisterConcrete(&FunnyAcct{}, "our/funnyacct", nil)

	// Setup the app
	app := mentapp.NewMockApp(sdk.DefaultJsonTxDecoder(cdc)) // inmemory tree

	// Set up initial chain state.  Bob has an account, alice doesn't yet
	app.OnInitialStart(func(ctx sdk.Context, req abci.RequestInitChain) (resp abci.ResponseInitChain) {
		// Create an account for bob in storage
		acct1 := &FunnyAcct{
			Owner:   bob.GetAddress(),
			Balance: uint32(10),
			PubKey:  bob.GetPubKey(),
		}
		bits, err := acct1.Encode()
		if err != nil {
			panic(err)
		}
		ctx.Db.Set(acct1.Owner, bits)
		return
	})

	// Add a validator
	app.OnVerifyTx(MoneyExchangeCheckTx)

	// Add a BeginBlock handler
	app.OnBeginBlock(func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		return abci.ResponseBeginBlock{
			Tags: sdk.Tags{
				sdk.Tag{Key: []byte("begin"), Value: []byte("av")},
			},
		}
	})

	app.Route(FunnyMoneyRouteKey, MoneyExchange)
	app.RouteQuery("/hello/account", QueryAccount)

	// Add an EndBlock handler
	app.OnEndBlock(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return abci.ResponseEndBlock{
			Tags: sdk.Tags{
				sdk.Tag{Key: []byte("end"), Value: []byte("av")},
			},
		}
	})
	return app
}
