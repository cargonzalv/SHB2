package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigWithValidValues(t *testing.T) {
	cfg, _ := NewConfig("config.yml")
	assert.NotNil(t, cfg, "cfg created with valid values")
}

func TestNewConfigWithEmptyValues(t *testing.T) {
	cfg, _ := NewConfig("")
	assert.NotNil(t, cfg, "cfg created with default values")
}

func TestNewConfigWithInvalidValues(t *testing.T) {
	cfg, _ := NewConfig("abc.yml")
	assert.Nil(t, cfg, "cfg could not be created")
}

func Test_load_env(t *testing.T) {
	cfg, _ := NewConfig("config.yml")
	loadEnv(cfg)
	assert.NotNil(t, cfg, "cfg created with valid values")
	assert.Equal(t, cfg.App.Environment, "dev")
}
