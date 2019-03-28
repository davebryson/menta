package test

import (
	"fmt"
	"testing"

	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
)

/*
 * Produces a message like this:
 {
 "type": "menta/stdtx",
 "value": {
  "msg": {
   "type": "our/hello",
   "value": {
    "name": "dave"
   }
  },
  "signer": "",
  "signature": null
 }
}
*/

func TestMsgEncoding(t *testing.T) {
	assert := assert.New(t)

	tx := sdk.StdTx{
		Msg: HelloMsg{Name: "dave"},
	}

	r, e := cdc.MarshalJSONIndent(tx, "", " ")
	assert.Nil(e)

	fmt.Printf("%v", string(r))

	txback := new(sdk.StdTx)
	err := cdc.UnmarshalJSON(r, txback)
	assert.Nil(err)
	assert.NotNil(txback.GetMsg())
	assert.Equal("hello", txback.GetMsg().Route())
	assert.Equal("info", txback.GetMsg().Type())

	hmsg := tx.GetMsg().(HelloMsg)
	assert.Equal("dave", hmsg.Name)
}
