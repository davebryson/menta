package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pv "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
)

// CreateNode creates an embedded tendermint node for standalone mode
func (app *MentaApp) CreateNode() *node.Node {
	// Assumes priv validator has been generated.  See setup()
	nodeKey, err := p2p.LoadOrGenNodeKey(app.Config.NodeKeyFile())
	if err != nil {
		panic(err)
	}

	node, err := node.NewNode(
		app.Config,
		pv.LoadOrGenFilePV(app.Config.PrivValidatorKeyFile(), app.Config.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(app.Config),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(app.Config.Instrumentation),
		logger,
	)
	if err != nil {
		panic(err)
	}
	return node
}

// Run run a standalone / in-process tendermint app
func (app *MentaApp) Run() {
	node := app.CreateNode()
	err := node.Start()
	if err != nil {
		logger.Error(err.Error())
	}
	TrapSignal(func() {
		if node.IsRunning() {
			node.Stop()
		}
	})
	select {}
}

// TrapSignal Adapted from Cosmos SDK
// Note: Must add a select{} after this - see above
func TrapSignal(cleanupFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		if cleanupFunc != nil {
			cleanupFunc()
		}
		exitCode := 128
		switch sig {
		case syscall.SIGINT:
			exitCode += int(syscall.SIGINT)
		case syscall.SIGTERM:
			exitCode += int(syscall.SIGTERM)
		}
		os.Exit(exitCode)
	}()
}
