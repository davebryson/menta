package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	menta "github.com/davebryson/menta/app"
	"github.com/davebryson/menta/examples/services/counter"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/urfave/cli"
)

const (
	HomeDir     = "counterdata"
	AliceWallet = "alice_wallet"
	rpcAddr     = "tcp://localhost:26657"
)

// This is the counter application
func createApp() *menta.MentaApp {
	// runs tendermint init if needed
	menta.InitTendermint(HomeDir)
	// setup the app
	app := menta.NewApp("counter-example", HomeDir)
	// add the service
	app.AddService(counter.CounterService{})

	return app
}

// RunApp sets up the menta application and starts the node
func RunApp() {
	app := createApp()
	app.Run()
}

func QueryCounter() {
	alice := counter.WalletFromSeed(AliceWallet)
	client := rpcclient.NewHTTP(rpcAddr, "/websocket")
	result, err := client.ABCIQuery(counter.ServiceName, alice.PubKey())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	count, err := counter.DecodeCount(result.Response.Value)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf(" ~~ Counter State for %x ~~\n", alice.PubKey())
	fmt.Printf(" ==> %v\n", count)
}

func SendTransaction(val uint32) {
	alice := counter.WalletFromSeed(AliceWallet)
	txbits, err := alice.NewTx(val)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := rpcclient.NewHTTP(rpcAddr, "/websocket")
	result, err := client.BroadcastTxCommit(txbits)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(resp))
}

// Simple command line application.  You should use two terminals for this:
// 1 to execute 'start' to run the tendermint application, and
// 1 to run the client.
func main() {
	app := cli.NewApp()
	app.Name = "counter cli"
	app.Version = "1.0"
	app.Description = "Menta counter example"
	app.Author = "Dave Bryson"
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Start the tendermint node",
			Action: func(c *cli.Context) error {
				RunApp()
				return nil
			},
		},
		{
			Name:  "send",
			Usage: "Send a transaction",
			Action: func(c *cli.Context) error {
				val := c.Args().Get(0)
				i, e := strconv.ParseInt(val, 0, 64)
				if e != nil {
					fmt.Printf("Error: '%v' is not a valid number\n", val)
					return nil
				}
				SendTransaction(uint32(i))
				return nil
			},
		},
		{
			Name:  "state",
			Usage: "Check state",
			Action: func(c *cli.Context) error {
				QueryCounter()
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
