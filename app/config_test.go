package app

import (
	"os"
	"testing"
	"time"

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
	cfg, err := LoadConfig("./nothere")
	assert.NotNil(err)
	assert.Nil(cfg)

	InitTendermint(TestConfigDir)
	a := NewApp("ex", TestConfigDir)
	assert.Equal("tcp://127.0.0.1:26658", a.Config.ProxyApp)

	assert.Panics(func() { NewApp("bad", "./bad") })
}
