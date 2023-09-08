package udws3

import (
	"testing"

	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/config"
	"github.com/stretchr/testify/assert"
)

func TestFetchFilesFromUdwS3BucketReturnsEmptysliceWhenErr(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")

	logger := log.GetLogger()
	udws3Client := NewService(logger, cfg)
	results, _ := udws3Client.FetchGzFiles()

	assert.Equal(t, len(results) == 0, true)
}

func TestDownloadGzFileReturnsFalseAndErrWhenErr(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")

	logger := log.GetLogger()
	udws3Client := NewService(logger, cfg)
	completed, err := udws3Client.DownloadGzFile("key", "/tmp/test.json.gz")
	assert.Equal(t, completed, false)
	assert.Equal(t, true, err != nil)
}
