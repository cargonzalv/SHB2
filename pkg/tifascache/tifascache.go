// tifascache package.
package tifascache

import (
	"strconv"
	"time"

	"github.com/TwiN/gocache/v2"
	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/config"
)

// Tifascache parameters.
type Tifascache struct {
	//logger service
	logger log.Service
	//tifascache service
	tifasCache *gocache.Cache
	//config service
	config *config.Config
}

// This statement forcing the module to
// implements the Tifascache `Service` interface.
var _ Service = (*Tifascache)(nil)

// NewService is a constructor function which get logger implementation and
// config params arguments and return implementation of Tifascache interface.
func NewService(l log.Service, c *config.Config) Service {
	cache := gocache.NewCache().WithMaxSize(gocache.NoMaxSize).WithEvictionPolicy(gocache.LeastRecentlyUsed)
	return &Tifascache{
		logger:     l,
		tifasCache: cache,
		config:     c,
	}
}

// GetTifa checks if the tifa exists in the tifascache
func (tc *Tifascache) GetTifa(tifa string) bool {
	_, exists := tc.tifasCache.Get(tifa)
	return exists
}

// SetTifa adds the tifa into the tifascache
func (tc *Tifascache) SetTifa(key string, val string, ttl time.Duration) {
	if ttl == 0 {
		defaultTtl, _ := strconv.Atoi(tc.config.Cache.KeyTTL)
		tc.tifasCache.SetWithTTL(key, val, time.Duration(defaultTtl)*time.Hour)
		return
	}
	tc.tifasCache.SetWithTTL(key, val, ttl)
}

// IsLastLoadTsExpired checks if the last tifas dump in the s3 bucket was expired or not
func (tc *Tifascache) IsLastLoadTsExpired(newLoadTs string) bool {
	prevLoadTs, exists := tc.tifasCache.Get(tc.config.Cache.LoadTimestampKey)
	if exists && prevLoadTs != nil {
		tc.logger.Debug("checking IsLastLoadTsExpired",
			log.Metadata{"newLoadTimestamp": newLoadTs, "prevLoadTimestamp": prevLoadTs,
				"expired": newLoadTs != prevLoadTs})
		return newLoadTs != prevLoadTs
	}
	return true
}
