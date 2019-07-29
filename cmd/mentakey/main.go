package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/davebryson/menta/crypto"
)

// Utility to create keys dumped to json

type key struct {
	Type       string
	Privatekey string
	Publickey  string
	Address    string
}

func makeKeys(num int) []byte {
	keys := make([]key, 0)
	for i := 0; i < num; i++ {
		sk := crypto.GeneratePrivateKey()
		v := key{
			Privatekey: sk.ToHex(),
			Publickey:  sk.PubKey().ToHex(),
		}
		keys = append(keys, v)
	}
	bits, err := json.MarshalIndent(keys, "", " ")
	if err != nil {
		panic(err)
	}
	return bits
}

// Help
var usage = ` Generate 1 or more cryptographic keys

 Usage: mentakey -num [val] 
 where [val] is the number of keys to create.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(0)
	}
	numbPtr := flag.Int("num", 1, "the number of keys to create")
	flag.Parse()

	keys := makeKeys(*numbPtr)
	fmt.Println(string(keys))
}
