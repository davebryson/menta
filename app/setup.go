package app

import (
	"os"
	"path/filepath"

	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
	pv "github.com/tendermint/tendermint/privval"
	tmtypes "github.com/tendermint/tendermint/types"
)

var chainIdPrefix = "menta-chain-%v"

func InitTendermint(homedir string) {
	if homedir == "" {
		homedir = DefaultHomeDir
	}
	if !cmn.FileExists(filepath.Join(homedir, "config", "config.toml")) {
		createConfig(homedir)
	}
}

func createConfig(homedir string) {
	config := cfg.DefaultConfig()
	if homedir == "" {
		config.SetRoot(DefaultHomeDir)
	} else {
		config.SetRoot(homedir)
	}

	cfg.EnsureRoot(config.RootDir)
	privValFile := config.PrivValidatorFile()
	privValidator := pv.LoadOrGenFilePV(privValFile)
	privValidator.Save()

	genFile := config.GenesisFile()
	chain_id := cmn.Fmt(chainIdPrefix, cmn.RandStr(6))

	// Create and save the genesis if it doesn't exist
	if _, err := os.Stat(genFile); os.IsNotExist(err) {
		// Set the chainid
		genDoc := tmtypes.GenesisDoc{ChainID: chain_id}
		// Add the validators
		genDoc.Validators = []tmtypes.GenesisValidator{tmtypes.GenesisValidator{
			PubKey: privValidator.PubKey,
			Power:  10,
		}}
		genDoc.SaveAs(genFile)
	}
}
