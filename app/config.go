package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
)

const (
	// MENTAHOME is the default dir for tendermint configuration
	MENTAHOME = ".menta"
	// Home is for viper configuration
	Home = "home"
)

// DefaultHomeDir for tendermint config
var DefaultHomeDir = os.ExpandEnv(fmt.Sprintf("$HOME/%s", MENTAHOME))

// LoadConfig using tendermint config
func LoadConfig(homedir string) (*cfg.Config, error) {
	if homedir == "" {
		homedir = DefaultHomeDir
	}

	if !tmos.FileExists(filepath.Join(homedir, "config", "config.toml")) {
		return nil, fmt.Errorf("Missing homedir! Did you run the init command?")
	}

	// Have a config file, load it
	viper.Set(Home, homedir)
	viper.SetConfigName("config")
	viper.AddConfigPath(homedir)
	viper.AddConfigPath(filepath.Join(homedir, "config"))

	// I don't think this ever returns an err.  It seems to create a default config if missing
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Missing homedir/config file. Did you run the init command?")
	}

	config := cfg.DefaultConfig()
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	config.SetRoot(config.RootDir)
	cfg.EnsureRoot(config.RootDir)

	return config, nil
}
