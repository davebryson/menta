# MENTA
A *simple* framework for creating Tendermint (permissioned) blockchain applications. 

But what about [Cosmos-SDK](https://github.com/cosmos/cosmos-sdk) ?.  Good question. If you're building a public application or planning to deploy to Cosmos, you should definitely use it. 

Menta is designed primarily for enterprise/permissioned blockchains where there's a need to increase the confidence in transactions among parties to the systems.

This framework is for:
* Enterprise blockchains that desire Byzantine Fault Tolerance 
* Rapid prototyping and small pilot projects  
* Folks looking to build Tendermint applications not destined for the Cosmos, or 
* Folks just wanting to learn *how* a Tendermint ABCI works

Menta provides a simple/minimal API on top of Tendermint based on our experience building many Tendermint applications from scratch. It also adopts some of the code and practices from the Cosmos SDK.

Of course, you can always start here and port to the Cosmos SDK later. That's the magic of Tendermint ABCI!

## Transaction Codec
Menta uses protobuf for the base transaction model. It's a minimal model so it's up to the user to decide how to encode/decode application specific messages. The Tx *wrapper* provides a way to route 
and transport application specific messages to menta handlers.

```
 message Tx {
   string route = 1;
   bytes msg = 3;
   bytes sender = 4;
   bytes nonce = 5;
   bytes sig = 6;
 }
```

* **route** is use to route transactions to a specific handler. It can also be used to help the application determine how to decode the `msg` payload.
* **msg** is an encoded application specific message.  How you encode the msg is up to you.
* **sender** is an optional field to store the wallet address of the sender
* **nonce** is an optional field to store a unique transaction nonce. Often used when signing the transaction
* **sig** is an optional field to store a cryptographic signature

`tx.go` in `types` provides functionality for signing and verifying transactions.

## Setup
**Current supported Tendermint version: v0.31.7**

Requires Go >= 1.12

Get menta: `go get github.com/davebryson/menta`

## Example
See `examples/counter` for a complete example of a simple application
