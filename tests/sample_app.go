package test

import (
	"encoding/json"

	mentapp "github.com/davebryson/menta/app"
	"github.com/davebryson/menta/codec"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	auth "github.com/davebryson/menta/plugins/authenticate"
)

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

func (wallet FunnyMoneyWallet) SendMoney(cdc *codec.Codec, msg auth.SendCoinMsg) ([]byte, error) {
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

// App
func createApp(bob, alice FunnyMoneyWallet) *mentapp.MentaApp {
	// Setup the app
	app := mentapp.NewMockApp(sdk.DefaultJSONTxDecoder(cdc)) // inmemory tree
	accountManager := auth.NewAccountManager(cdc)

	// Set up initial chain state.
	app.OnInitialStart(func(ctx sdk.Context, req abci.RequestInitChain) (resp abci.ResponseInitChain) {
		// Create accounts for bob & alice in storage
		acct := &auth.BasicAccount{
			Owner:   bob.GetAddress(),
			Balance: uint32(10),
			PubKey:  bob.GetPubKey(),
		}
		aliceacct := &auth.BasicAccount{
			Owner:   alice.GetAddress(),
			Balance: uint32(0),
			PubKey:  alice.GetPubKey(),
		}
		err := accountManager.SaveAccount(ctx, acct)
		if err != nil {
			panic(err)
		}
		err = accountManager.SaveAccount(ctx, aliceacct)
		if err != nil {
			panic(err)
		}

		return
	})

	// Add a validator
	app.OnVerifyTx(auth.VerifyAccountHandler(accountManager))

	// Add a BeginBlock handler
	app.OnBeginBlock(func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		return abci.ResponseBeginBlock{
			Tags: sdk.Tags{
				sdk.Tag{Key: []byte("begin"), Value: []byte("av")},
			},
		}
	})

	app.Route(auth.AuthenticateRoute, auth.AccountTxHandler(accountManager))
	app.RouteQuery(auth.AuthenticateQueryRoute, auth.AccountQueryHandler(accountManager))

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
