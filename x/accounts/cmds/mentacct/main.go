package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/davebryson/menta/x/accounts"
	"github.com/urfave/cli"
)

// Generate 'num' private keys to JSON
// Example: mentacct 2 > accounts.json
func main() {
	app := cli.NewApp()
	app.Name = "mentacct"
	app.Usage = "mentacct [num]"
	app.Description = "generate keys in JSON"
	app.Action = func(c *cli.Context) error {
		val, err := strconv.Atoi(c.Args().Get(0))
		if err != nil {
			fmt.Printf("Error: '%v' is not a valid number\n", val)
			return nil
		}
		result, err := accounts.GenerateJSONAccounts(val)
		if err != nil {
			fmt.Println("Error generating accounts")
			return nil
		}
		fmt.Printf("%s", result)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
