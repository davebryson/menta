package app

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/types"
)

const TestDir = "test_app_config"

func TestNodeApp(t *testing.T) {
	defer func() {
		time.Sleep(2 * time.Second)
		os.RemoveAll(TestDir)
	}()

	InitTendermint(TestDir)
	app := NewApp("testapp", TestDir)

	node := app.CreateNode()
	err := node.Start()
	if err != nil {
		t.Error(err)
	}

	// Adapted from Tendermint node_test...
	blockSub, err := node.EventBus().Subscribe(context.Background(), "node_app_test", types.EventQueryNewBlock)
	assert.NoError(t, err)
	select {
	case <-blockSub.Out():
	case <-blockSub.Cancelled():
		t.Fatal("blocksSub was cancelled")
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
