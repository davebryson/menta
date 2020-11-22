package client

// Small client API to connect to Tendermint

import (
	"context"
	"encoding/json"

	sdk "github.com/davebryson/menta/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client/http"
)

const rpcAddr = "tcp://localhost:26657"

// SendTx : Send a signed transaction to a local node.
// Useful for local command line clients
func SendTx(tx *sdk.SignedTransaction) (string, error) {
	encodedMsg, err := sdk.EncodeTx(tx)
	if err != nil {
		return "", err
	}

	client, _ := rpcclient.New(rpcAddr, "/websocket")
	result, err := client.BroadcastTxCommit(context.Background(), encodedMsg)
	if err != nil {
		return "", err
	}

	resp, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// Query the state of a given service
func Query(serviceName string, key []byte) ([]byte, error) {
	client, _ := rpcclient.New(rpcAddr, "/websocket")
	result, err := client.ABCIQuery(context.Background(), serviceName, key)
	if err != nil {
		return nil, err
	}
	return result.Response.Value, nil
}
