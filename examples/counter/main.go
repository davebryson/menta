package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	counterapp "github.com/davebryson/menta/examples/counter/app"
	"github.com/urfave/cli"
)

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
				counterapp.RunApp()
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
				counterapp.SendTx(uint32(i))
				return nil
			},
		},
		{
			Name:  "state",
			Usage: "Check state",
			Action: func(c *cli.Context) error {
				counterapp.CheckState()
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
