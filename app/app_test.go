package app

import (
	"testing"

	"github.com/davebryson/menta/examples/services/counter"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func createApp() *MentaApp {
	app := NewMockApp() // inmemory tree
	app.AddService(&counter.Service{})
	return app
}

// Test to check all callbacks and handler hooks
func TestAppCallbacks(t *testing.T) {
	assert := assert.New(t)
	app := createApp()

	alice := counter.CreateWallet()

	// --- Simulate running it ---

	// Call InitChain
	icresult := app.InitChain(abci.RequestInitChain{})
	assert.Equal(icresult, abci.ResponseInitChain{})
	// Commit here so we have something for Info
	c1 := app.Commit()

	// Call Info
	result := app.Info(abci.RequestInfo{})
	assert.Equal("mockapp", result.GetData())
	// block height should be 1 because we committed
	assert.Equal(int64(1), result.GetLastBlockHeight())
	hash1 := result.GetLastBlockAppHash()
	assert.NotNil(hash1)
	// Should == the first commit hash
	assert.Equal(c1.Data, hash1)

	// Fail: Call Query
	respQ := app.Query(abci.RequestQuery{Path: counter.ServiceName, Data: alice.PubKey()})
	assert.Equal(uint32(1), respQ.Code)

	// Pass: Run checkTx
	tx, err := alice.NewTx(1)
	assert.Nil(err)
	chtx := app.CheckTx(abci.RequestCheckTx{Tx: tx})
	assert.Equal(abci.ResponseCheckTx{Code: 0}, chtx)

	// Run Deliver handlers
	dtx := app.DeliverTx(abci.RequestDeliverTx{Tx: tx})
	assert.Equal(uint32(0), dtx.Code)

	// Commit the new state to storage
	commit := app.Commit()

	assert.NotNil(commit.Data)
	// Should be a new apphash
	assert.NotEqual(c1.Data, commit.Data)

	// Now committed state should == 1
	respQ = app.Query(abci.RequestQuery{Path: counter.ServiceName, Data: alice.PubKey()})
	assert.Equal(uint32(0), respQ.Code)

	count, err := counter.DecodeCount(respQ.GetValue())
	assert.Nil(err)
	assert.Equal(uint32(1), count.Current)
}
