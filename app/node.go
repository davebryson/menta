package app

import (
	"os"

	"github.com/tendermint/tendermint/abci/server"
	cmn "github.com/tendermint/tendermint/libs/common"
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
	node.Start()

	cmn.TrapSignal(func() {
		node.Stop()
	})
}

// RunServer starts a separate server that connects to tendermint
func (app *MentaApp) RunServer() {
	srv, err := server.NewServer("0.0.0.0:26658", "socket", app)
	if err != nil {
		cmn.Exit(err.Error())
	}
	srv.Start()

	cmn.TrapSignal(func() {
		srv.Stop()
	})
}
