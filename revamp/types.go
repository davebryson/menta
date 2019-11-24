package revamp

import (
	"github.com/davebryson/menta/store"
	sdk "github.com/davebryson/menta/types"
	amino "github.com/tendermint/go-amino"
)

var codec = amino.NewCodec()

// Message is part of the transaction
// Called by the Service
type Message interface {
	Execute(sender []byte, store sdk.KVStore) sdk.Result
}

type SignedTransaction struct {
	Route     string
	Sender    []byte
	Msg       Message
	Nonce     []byte
	Signature []byte
}

type Service interface {
	Route() string
	RegisterMessages(cdc *amino.Codec)
	OnGenesis(store sdk.KVStore)
	Execute(sender []byte, msg Message, store sdk.KVStore) sdk.Result
	Query(key []byte, store sdk.KVStore)
}

func RegisterMentaTypes(cdc *amino.Codec) {
	cdc.RegisterInterface((*Message)(nil), nil)
	cdc.RegisterConcrete(&SignedTransaction{}, "menta/signedtx", nil)
}

// Temp sim engine...
type SimApp struct {
	codec    *amino.Codec
	services map[string]Service
	dcache   *store.KVCache
}

func NewSimApp() *SimApp {
	st := store.NewStateStore("")
	dcache := store.NewCache(st)
	codec := amino.NewCodec()
	RegisterMentaTypes(codec)

	return &SimApp{
		codec:    codec,
		services: make(map[string]Service, 0),
		dcache:   dcache,
	}
}

func (sa *SimApp) AddService(s Service) {
	s.RegisterMessages(sa.codec)
	sa.services[s.Route()] = s
}

func (sa *SimApp) Run(raw []byte) sdk.Result {
	tx := new(SignedTransaction)
	sa.codec.UnmarshalBinaryLengthPrefixed(raw, tx)

	sender := tx.Sender
	msg := tx.Msg
	return msg.Execute(sender, sa.dcache)
}
