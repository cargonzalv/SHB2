package tifascache

import (
	"testing"
	"time"

	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/config"
	"github.com/stretchr/testify/assert"
)

func TestIfTifaExists(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")

	logger := log.GetLogger()
	cacheClient := NewService(logger, cfg)
	cacheClient.SetTifa("tifaKey", "tifaVal", 0)
	exists := cacheClient.GetTifa("tifaKey")

	assert.Equal(t, exists, true)
}

func TestIfTifaExistsAfterTtlExpires(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")

	logger := log.GetLogger()
	cacheClient := NewService(logger, cfg)
	cacheClient.SetTifa("tifaKey", "tifaVal", 1*time.Millisecond)
	time.Sleep(1 * time.Millisecond)
	exists := cacheClient.GetTifa("tifaKey")

	assert.Equal(t, exists, false)
}
