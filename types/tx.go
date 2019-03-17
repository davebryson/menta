package types

import (
	fmt "fmt"

	wire "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/common"
)

// *** Most of this code is adapted directly from the Cosmos SDK. ***

// Msg is the application specific message(s) that invoke handlers
type Msg interface {
	// Return the message type.
	Route() string

	// Returns a human-readable string for the message
	Type() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error // TODO make work with result
}

// Tx - Sender submits 1 msg and the associated signature
type Tx interface {
	GetMsg() Msg

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	GetSigner() common.HexBytes

	ValidateBasic() error
}

var _ Tx = (*StdTx)(nil)

type StdTx struct {
	Msg       Msg             `json:"msg"`
	Signer    common.HexBytes `json:"signer"`
	Signature []byte          `json:"signature"`
}

func (tx StdTx) GetMsg() Msg {
	return tx.Msg
}

func (tx StdTx) GetSignBytes() []byte {
	return tx.Signature
}

func (tx StdTx) GetSigner() common.HexBytes {
	return tx.Signer
}

func (tx StdTx) ValidateBasic() error {
	if tx.Msg == nil {
		return fmt.Errorf("no msgs in the Tx - nothing to do")
	}
	if tx.Signature == nil {
		return fmt.Errorf("missing signature")
	}
	if tx.Signer == nil {
		return fmt.Errorf("missing signer address")
	}
	return nil
}

type TxDecoder func(txBytes []byte) (Tx, error)

type TxEncoder func(tx Tx) ([]byte, error)

func DefaultJsonTxDecoder(cdc *wire.Codec) TxDecoder {
	return func(txBytes []byte) (Tx, error) {
		var tx = StdTx{}
		if len(txBytes) == 0 {
			return nil, fmt.Errorf("tx is empty")
		}
		err := cdc.UnmarshalJSON(txBytes, &tx)
		if err != nil {
			return nil, err
		}
		return tx, nil
	}
}
