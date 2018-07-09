package app

import (
	"os"

	"github.com/tendermint/tendermint/abci/server"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	pv "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
)

func (app *MentaApp) CreateNode() *node.Node {
	// Assumes priv validator has been generated.  See loadConfig()
	//privValFile := app.config.PrivValidatorFile()
	//privValidator := pv.LoadFilePV(privValFile)
	//papp := proxy.NewLocalClientCreator(app)
	node, err := node.NewNode(
		app.config,
		pv.LoadOrGenFilePV(app.config.PrivValidatorFile()),
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(app.config),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider,
		logger,
	)
	if err != nil {
		panic(err)
	}
	return node
}

func (app *MentaApp) Run() {
	node := app.CreateNode()
	node.Start()

	cmn.TrapSignal(func() {
		node.Stop()
	})
}

func (app *MentaApp) RunServer() {
	srv, err := server.NewServer("0.0.0.0:46658", "socket", app)
	if err != nil {
		cmn.Exit(err.Error())
	}
	srv.Start()

	cmn.TrapSignal(func() {
		srv.Stop()
	})
}
