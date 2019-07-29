package app

import (
	fmt "fmt"

	mentacli "github.com/davebryson/menta/client"
	sdk "github.com/davebryson/menta/types"
)

func makeTx(val uint32) *sdk.Tx {
	encoded := encodeCount(val)
	return &sdk.Tx{Route: routeName, Msg: encoded}
}

// SendTx is called from the cli to send a transaction to the node
func SendTx(val uint32) {
	msg := makeTx(val)
	resp, err := mentacli.SendTx(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(resp))
}

// CheckState queries the node for the latest confirmed state of the count
func CheckState() {
	result, err := mentacli.Query(queryRoute, stateKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf(" State is: %d\n", decodeCount(result))
}
