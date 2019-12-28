package testkit

import (
	"time"

	menta "github.com/davebryson/menta/app"
	sdk "github.com/davebryson/menta/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

const rpcAddr = "tcp://localhost:26657"

// TestKit provides a simple API to test your service against a running node.
type TestKit struct {
	app         *menta.MentaApp
	client      *rpcclient.HTTP
	serviceName string
	homedir     string
}

// NewTestKit creates an in-process Tendermint node.
// homedir: is the directory where Tendermint will store configuration.  You should use a
//          'defer os.RemoveAll()' at the start of your test to delete this directory
//           when the test finishes
// service: is you application service to test
func NewTestKit(homedir string, service sdk.Service) TestKit {
	menta.InitTendermint(homedir)
	app := menta.NewApp("testkit", homedir)
	app.AddService(service)
	return TestKit{
		app:         app,
		client:      rpcclient.NewHTTP(rpcAddr, "/websocket"),
		serviceName: service.Name(),
		homedir:     homedir,
	}
}

// Launch start the Tendermint node and waits for 2 seconds to ensure the node is running
func (tk TestKit) Launch() {
	go func() {
		tk.app.Run()
	}()
	// Wait to make sure all is running
	time.Sleep(2 * time.Second)
}

// Query your service for the given key. The key should match what's
// expected when you implemented service.Query(...)
func (tk TestKit) Query(key []byte) (*core_types.ResultABCIQuery, error) {
	return tk.client.ABCIQuery(tk.serviceName, key)
}

// SendTxCommit sends a transaction and waits for it to be committed
func (tk TestKit) SendTxCommit(txencoded []byte) (*core_types.ResultBroadcastTxCommit, error) {
	return tk.client.BroadcastTxCommit(txencoded)
}

// SendTxAsync can be used to send several transactions at once.  It doesn't wait for commit
func (tk TestKit) SendTxAsync(txencoded []byte) (*core_types.ResultBroadcastTx, error) {
	return tk.client.BroadcastTxAsync(txencoded)
}
