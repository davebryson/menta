package types

// Service is the primary interface to implement for application services.
// A given MentaApp may have 1 or more of these.
type Service interface {
	// Route is the unique name of the service. Used to register your service in Menta
	Route() string
	// Init is called once, on the very first run of the application.
	// Use this to load genesis data for your service
	Init(TxContext)
	// Execute is the primary business logic of your service. This is the blockchain
	// state transistion function
	Execute(TxContext) Result
	// Query provides read access to service storage.
	Query([]byte, QueryContext) Result
}

// ValidateTxHandler should be implemented to validate/check a transaction for
// inclusion into the mempool.  This is called on 'checkTx'.  A returned non-zero
// result.Code will exclude a transaction from consideration.  A Menta application has
// only 1 of these
type ValidateTxHandler func(ctx TxContext) Result
