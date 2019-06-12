# Counter Example with Menta

This shows a simple example app built with Menta.  It's the common 'counter' application.

The application starts with an initial state of zero.  The client sends a number to the app to 
increment the state.  But, they must send the next expected value. For example, if the current
application state = 1. The client must send 2, otherwise the transaction is rejected.

## Setup
Note: require go version 1.12 or greater.



