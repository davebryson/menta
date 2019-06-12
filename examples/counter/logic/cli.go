package logic

import (
	"encoding/json"
	fmt "fmt"

	sdk "github.com/davebryson/menta/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

const rpcAddr = "tcp://localhost:26657"

func makeTx(val uint32) ([]byte, error) {
	encoded := encodeCount(val)
	t := &sdk.Tx{Route: routeName, Msg: encoded}
	return sdk.EncodeTx(t)
}

/*func makeTx(val int64) ([]byte, error) {
	// Encode the application specific message
	msg := &CountMsg{Value: val}
	encodedMsg, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	// Create a basic menta transaction
	t := &sdk.Tx{
		Route: routeName,
		Msg:   encodedMsg,
	}
	// Encode it for transport
	return sdk.EncodeTx(t)
}*/

// SendTx is called from the cli to send a transaction to the node
func SendTx(val uint32) {
	msg, err := makeTx(val)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := rpcclient.NewHTTP(rpcAddr, "/websocket")
	result, err := client.BroadcastTxCommit(msg)
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

// CheckState queries the node for the latest confirmed state of the count
func CheckState() {
	client := rpcclient.NewHTTP(rpcAddr, "/websocket")
	result, err := client.ABCIQuery(queryRoute, stateKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf(" State is: %d\n", decodeCount(result.Response.Value))
}
