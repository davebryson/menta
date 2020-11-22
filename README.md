# MENTA
A *lighweight* framework for creating Tendermint (permissioned) blockchain applications. 

Menta is designed primarily for enterprise/permissioned blockchains where there's a need to increase confidence in transactions among parties to the systems.

This framework is for:
* Permissioned blockchains that desire Byzantine Fault Tolerance 
* Rapid prototyping and small pilot projects
* Folks looking to build Tendermint applications not destined for the Cosmos, or 
* Folks just wanting to learn *how* a Tendermint ABCI works

Menta provides a simple API on top of Tendermint based on our experience building many Tendermint applications from scratch. It also adopts some of the code and practices from the Cosmos SDK.

Why not use [Cosmos-SDK](https://github.com/cosmos/cosmos-sdk)?  Good question. If you're building a public application or planning to deploy to Cosmos, you should definitely use it. 

Of course, you can always start here and port to the Cosmos SDK later. That's the magic of Tendermint ABCI!

## Transaction Codec
Menta uses protobuf for the base transaction model. It's a minimal model leaving it up to the user to decide how to encode/decode application specific messages (msg field). The Tx *wrapper* provides a way to route and transport application specific messages to menta handlers.

```
 message Tx {
   string service = 1;
   bytes msg = 2;
   uint32 msgid = 3
   bytes sender = 4;
   bytes nonce = 5;
   bytes sig = 6;
 }
```

* **service** is use to route transactions to a specific `Service`.
* **msg** is an encoded application specific message.  How you encode the msg is up to you.
* **msgid** can be used to distinquish messages for decoding
* **sender** is an optional field to store the wallet address of the sender
* **nonce** is an optional field to store a unique transaction nonce. Often used when signing the transaction
* **sig** is an optional field to store a cryptographic signature

`tx.go` in `types` provides functionality for signing and verifying transactions.

## Setup
**Current supported Tendermint version: v0.34.0**

Requires Go >= 1.15

Get menta: `go get -u github.com/davebryson/menta`

## Example
See `examples/counter` for a complete example of a simple application
