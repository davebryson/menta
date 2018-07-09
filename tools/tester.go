package tools

import (
	"github.com/davebryson/menta/app"
	sdk "github.com/davebryson/menta/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Tester - simple tool to test tx logic without tendermint
// Handles initChain, CheckTx, DeliverTx, Commit.
// It currently does not do Begin/EndBlock
type Tester struct {
	app            *app.MentaApp
	commit         sdk.CommitInfo
	mempool        [][]byte
	consensuspool  [][]byte
	numDeliveredTx int
}

// NewTester - setup the app.  First create a new app via app.NewTestApp()
// Attach callbacks, then pass to this function
func NewTester(app *app.MentaApp) *Tester {
	t := &Tester{
		app:           app,
		mempool:       make([][]byte, 0),
		consensuspool: make([][]byte, 0),
	}
	t.app.InitChain(abci.RequestInitChain{})
	t.app.Commit()
	return t
}

// Call one or more times
func (tst *Tester) SendTx(tx []byte) {
	tst.mempool = append(tst.mempool, tx)
}

// Runs the mempool and commits data
// Returns CommitInfo and the num of txs actually processed (delivered)
// This can vary based on the result of CheckTx validation
func (tst *Tester) MakeBlock() (sdk.CommitInfo, int) {
	// CheckTx
	for _, tx := range tst.mempool {
		// could measure tx size in bytes
		resp := tst.app.CheckTx(tx)
		if resp.Code == 0 {
			tst.consensuspool = append(tst.consensuspool, tx)
		}
	}

	tst.app.BeginBlock(abci.RequestBeginBlock{})

	// DeliverTx
	for _, tx := range tst.consensuspool {
		// Could measure execution time here
		resp := tst.app.DeliverTx(tx)
		if resp.Code == 0 {
			tst.numDeliveredTx += 1
		}
	}

	tst.app.EndBlock(abci.RequestEndBlock{})

	// Commit
	hash := tst.app.Commit().Data
	version := tst.commit.Version + 1
	tst.commit = sdk.CommitInfo{
		Hash:    hash,
		Version: version,
	}

	totalTx := tst.numDeliveredTx
	// reset pools/stats
	tst.mempool = make([][]byte, 0)
	tst.consensuspool = make([][]byte, 0)
	tst.numDeliveredTx = 0

	return tst.commit, totalTx
}

// Query *committed* state
func (tst *Tester) QueryByKey(key string) (uint32, []byte) {
	resp := tst.app.Query(abci.RequestQuery{Path: "/key", Data: []byte(key)})
	return resp.Code, resp.Value
}
