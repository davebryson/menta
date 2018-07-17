package test

import (
	"context"
	"os"
	"testing"
	"time"

	menta "github.com/davebryson/menta/app"
	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/types"
)

const TEST_DIR = "test_app_config"

func TestNodeApp(t *testing.T) {
	defer func() {
		time.Sleep(1 * time.Second)
		os.RemoveAll(TEST_DIR)
	}()

	menta.InitTendermint(TEST_DIR)
	app := menta.NewApp("testapp", TEST_DIR)
	app.OnTx("hello", func(ctx sdk.Context) sdk.Result {
		return sdk.Result{}
	})

	node := app.CreateNode()
	err := node.Start()
	if err != nil {
		t.Error(err)
	}

	// Adapted from Tendermint/node/node_test...
	blockCh := make(chan interface{})
	err = node.EventBus().Subscribe(context.Background(), "node_app_test", types.EventQueryNewBlock, blockCh)
	assert.NoError(t, err)
	select {
	case <-blockCh:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for the node to produce a block")
	}

	// stop the node
	go func() {
		node.Stop()
	}()

	select {
	case <-node.Quit():
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for shutdown")
	}
}
