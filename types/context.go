package types

// Context provides access to the state store and the current decoded Transaction
type Context struct {
	Db Cache
	Tx Tx
}

// NewContext creates a new context (hence the name)
func NewContext(db Cache, tx Tx) Context {
	return Context{
		Db: db,
		Tx: tx,
	}
}
