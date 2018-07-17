package test

import (
	"os"
	"testing"
	"time"

	"github.com/davebryson/menta/app"
	"github.com/stretchr/testify/assert"
)

const TEST_CONFIG_DIR = "./test_config"

func TestConfig(t *testing.T) {
	defer func() {
		time.Sleep(1 * time.Second)
		os.RemoveAll(TEST_CONFIG_DIR)
	}()

	assert := assert.New(t)
	cfg, err := app.LoadConfig("./nothere")
	assert.NotNil(err)
	assert.Nil(cfg)

	app.InitTendermint(TEST_CONFIG_DIR)
	a := app.NewApp("ex", TEST_CONFIG_DIR)
	assert.Equal("tcp://127.0.0.1:26658", a.Config.ProxyApp)

	assert.Panics(func() { app.NewApp("bad", "./bad") })
}
