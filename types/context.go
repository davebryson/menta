package types

import (
	fmt "fmt"
)

// Creates the store prefix key.  Storage is scoped by the
// service name
func createServiceKey(serviceName string, key []byte) []byte {
	return append([]byte(fmt.Sprintf("%s//", serviceName)), key...)
}

// TxContext provides access to storage and the current transaction.
// Store keys are prefixed by the name of the service.
// It's automatically passed to the Service
type TxContext struct {
	store       Cache
	Tx          *SignedTransaction
	serviceName string
}

// NewContext create a new instance. This is automatically called from Menta
func NewContext(serviceName string, store Cache, tx *SignedTransaction) TxContext {
	return TxContext{
		store:       store,
		Tx:          tx,
		serviceName: serviceName,
	}
}

// Insert a k/v into storage
func (sc TxContext) Insert(key, value []byte) {
	prefixedKey := createServiceKey(sc.serviceName, key)
	sc.store.Set(prefixedKey, value)
}

// Delete a k/v from storage
func (sc TxContext) Delete(key []byte) {
	prefixedKey := createServiceKey(sc.serviceName, key)
	sc.store.Delete(prefixedKey)
}

// Get a value from the current cache or storage
func (sc TxContext) Get(key []byte) ([]byte, error) {
	prefixedKey := createServiceKey(sc.serviceName, key)
	return sc.store.Get(prefixedKey)
}

// QueryContext provides read-only access to storage
type QueryContext struct {
	store       Cache
	serviceName string
}

// NewQueryContext create a new instance. This is called
// automatically by Menta
func NewQueryContext(serviceName string, store Cache) QueryContext {
	return QueryContext{
		store:       store,
		serviceName: serviceName,
	}
}

// Get returns a value for a given key from committed data
func (sc QueryContext) Get(key []byte) ([]byte, error) {
	prefixedKey := createServiceKey(sc.serviceName, key)
	return sc.store.GetCommitted(prefixedKey)
}
