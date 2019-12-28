package testkit

import (
	"os"
	"testing"
	"time"

	"github.com/davebryson/menta/examples/services/counter"
	"github.com/stretchr/testify/assert"
)

// director for tendermint configuration
const testdir = "./sample"

func TestApi(t *testing.T) {
	assert := assert.New(t)
	// Make sure to remove after testing
	defer func() {
		os.RemoveAll(testdir)
	}()

	// Set up the service
	service := counter.Service{}
	serviceName := service.Name()

	// Create the testkey with your service
	tester := NewTestKit(testdir, counter.Service{})
	// Start the node
	tester.Launch()

	// Create a wallet to send txs
	alice := WalletFromSeed("//Alice")
	msg := &counter.Increment{Value: 1}
	// Create tx
	txbits, err := alice.CreateTx(serviceName, 0, msg)
	assert.Nil(err)
	// Send the tx
	result, err := tester.SendTxCommit(txbits)
	assert.Nil(err)
	assert.Equal(uint32(0), result.DeliverTx.GetCode())

	// Send a batch of txs
	for i := 2; i < 1000; i++ {
		msg := &counter.Increment{Value: uint32(i)}
		txbits, err := alice.CreateTx(serviceName, 0, msg)
		assert.Nil(err)
		_, err = tester.SendTxAsync(txbits)
		assert.Nil(err)
		assert.Equal(uint32(0), result.DeliverTx.GetCode())
	}

	// Wait 2 blocks to allow commit since we used ..async
	time.Sleep(2 * time.Second)

	// Check committed state is correct
	// We store counters keyed by the senders publickey
	pubkey := alice.PubKey()
	// Query for the count
	rq, err := tester.Query(pubkey)
	assert.Nil(err)
	assert.Equal(uint32(0), rq.Response.GetCode())
	cv, err := counter.DecodeCount(rq.Response.GetValue())
	assert.Nil(err)
	assert.Equal(uint32(999), cv.Current)
}
