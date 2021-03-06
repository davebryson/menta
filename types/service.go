package types

// Service is the primary interface to implement for application services.
// A given MentaApp may have 1 or more of these.
type Service interface {
	// Name is the unique name of the service. Used to register your service in Menta
	Name() string
	// Init is called once, on the very first run of the application.
	// Use this to load genesis data for your service
	Initialize(data []byte, store Cache)
	// Execute is the primary business logic of your service. This is the blockchain
	// state transistion function
	Execute(sender []byte, msgid uint32, message []byte, store Cache) Result
	// Query provides read access to storage.
	Query(key []byte, store Snapshot) Result
}
