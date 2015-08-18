package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/directory"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) TestNameAndPaths() {
	cfg := new(viperProvider)
	assert := assert.New(suite.T())

	cfg.initNameAndPaths()
	assert.Equal("config", cfg.name)
	assert.Equal(defaultPaths, cfg.paths)

	cfg.initNameAndPaths("wercker-config", "/etc/foodtastechess")
	assert.Equal("wercker-config", cfg.name)
	assert.Equal([]string{"/etc/foodtastechess"}, cfg.paths)
}

func (suite *ConfigTestSuite) TestLogging() {
	cfg := NewConfigProvider("testconfig", "../").(*viperProvider)
	assert := assert.New(suite.T())

	assert.NotNil(cfg.log)
}

type cacheService struct {
	Config DatabaseConfig `inject:"databaseConfig"`
}

func (suite *ConfigTestSuite) TestProvide() {
	var (
		d directory.Directory = directory.New()
	)

	cfg := NewConfigProvider("testconfig", "../").(*viperProvider)

	d.AddService("configProvider", cfg)
	d.AddService("cacheService", &cacheService{})

	err := d.Start()

	assert := assert.New(suite.T())
	assert.Nil(err, "Unable to populate queries cache config")
}

// Entrypoint
func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
