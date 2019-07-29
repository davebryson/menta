package client

// Small client API to connect to Tendermint

import (
	"encoding/json"

	sdk "github.com/davebryson/menta/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

const rpcAddr = "tcp://localhost:26657"

// SendTx to a local node. Useful for local commandline clients
func SendTx(tx *sdk.Tx) (string, error) {
	encodedMsg, err := sdk.EncodeTx(tx)
	if err != nil {
		return "", err
	}

	client := rpcclient.NewHTTP(rpcAddr, "/websocket")
	result, err := client.BroadcastTxCommit(encodedMsg)
	if err != nil {
		return "", err
	}

	resp, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// Query a state route
func Query(route string, key []byte) ([]byte, error) {
	client := rpcclient.NewHTTP(rpcAddr, "/websocket")
	result, err := client.ABCIQuery(route, key)
	if err != nil {
		return nil, err
	}
	return result.Response.Value, nil
}
