package types

import "errors"

type TxContext struct {
	Service string
	Sender  []byte
	MsgID   uint32
	Message []byte
	rwstore KVStore
}

func NewTxContext(service string, sender []byte, msgid uint32, message []byte, store KVStore) TxContext {
	return TxContext{
		Service: service,
		Sender:  sender,
		MsgID:   msgid,
		Message: message,
		rwstore: store,
	}
}

// Store returns the k/v store for the current service
func (ctx TxContext) Store() PrefixedRWStore {
	return NewPrefixRW(ctx.Service, ctx.rwstore)
}

// StoreForService returns the k/v store for a another service with the given name.
// Returns an error if the service 'name' is not a registered service
func (ctx TxContext) StoreForService(name string) (PrefixedRWStore, error) {
	if ctx.rwstore.Has(prefixKey([]byte(REGISTERED_SERVICE_PREFIX), []byte(name))) {
		return PrefixedRWStore{}, errors.New("Service not found")
	}
	return NewPrefixRW(name, ctx.rwstore), nil
}
