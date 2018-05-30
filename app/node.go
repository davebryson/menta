package app

import (
	"os"

	"github.com/tendermint/abci/server"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	pv "github.com/tendermint/tendermint/types/priv_validator"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
)

func (app *MentaApp) CreateNode() *node.Node {
	// Assumes priv validator has been generated.  See loadConfig()
	privValFile := app.config.PrivValidatorFile()
	privValidator := pv.LoadFilePV(privValFile)
	papp := proxy.NewLocalClientCreator(app)
	node, err := node.NewNode(
		app.config,
		privValidator,
		papp,
		node.DefaultGenesisDocProviderFunc(app.config),
		node.DefaultDBProvider, logger,
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
