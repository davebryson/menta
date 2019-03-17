package test

import (
	"os"
	"testing"
	"time"

	"github.com/davebryson/menta/app"
	"github.com/stretchr/testify/assert"
)

const TestConfigDir = "./test_config"

// Test LoadConfig and InitTendermint
func TestConfig(t *testing.T) {
	defer func() {
		time.Sleep(1 * time.Second)
		os.RemoveAll(TestConfigDir)
	}()

	assert := assert.New(t)
	cfg, err := app.LoadConfig("./nothere")
	assert.NotNil(err)
	assert.Nil(cfg)

	app.InitTendermint(TestConfigDir)
	a := app.NewApp("ex", TestConfigDir, nil)
	assert.Equal("tcp://127.0.0.1:26658", a.Config.ProxyApp)

	assert.Panics(func() { app.NewApp("bad", "./bad", nil) })
}
