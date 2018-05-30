package app

import (
	"fmt"
	"os"

	cfg "github.com/tendermint/tendermint/config"
)

const (
	MENTAHOME = ".menta"
)

var DefaultHomeDir = os.ExpandEnv(fmt.Sprintf("$HOME/%s", MENTAHOME))

func LoadConfig(homedir string) (*cfg.Config, error) {
	config := cfg.DefaultConfig()
	if homedir == "" {
		config.SetRoot(DefaultHomeDir)
	} else {
		config.SetRoot(homedir)
	}
	cfg.EnsureRoot(config.RootDir)

	if _, err := os.Stat(config.PrivValidatorFile()); os.IsNotExist(err) {
		return config, fmt.Errorf("Missing homedir! Did you run the init command?")
	}
	return config, nil
}
