package types

type Context struct {
	Db Cache
	Tx *Transaction
}

func NewContext(db Cache, tx *Transaction) Context {
	return Context{
		Db: db,
		Tx: tx,
	}
}
