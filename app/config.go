package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	// MentaHome is the default home directory for the app
	MentaHome = ".menta"
	// Home is the configuration setting name for the home dir
	Home = "home"
)

// DefaultHomeDir is just that
var DefaultHomeDir = os.ExpandEnv(fmt.Sprintf("$HOME/%s", MentaHome))

// LoadConfig from the given dir
func LoadConfig(homedir string) (*cfg.Config, error) {
	if homedir == "" {
		homedir = DefaultHomeDir
	}

	if !cmn.FileExists(filepath.Join(homedir, "config", "config.toml")) {
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
