package mockservices

import (
	"github.com/adgear/go-commons/pkg/httpclient"
	"github.com/adgear/go-commons/pkg/kafka/producer"
	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/pkg/demand"
	"github.com/adgear/sps-header-bidder/pkg/tifascache"
	"github.com/adgear/sps-header-bidder/pkg/udws3"
)

//go:generate mockgen -destination=mocks.gen.go -source=dependencies.go -package=mockservices
type (
	logger interface {
		log.Service
	}

	demandClient interface {
		demand.Service
	}

	httpClient interface {
		httpclient.Service
	}

	kafkaProducer interface {
		producer.Service
	}

	cacheClient interface {
		tifascache.Service
	}

	udws3Client interface {
		udws3.Service
	}
)

var _ logger = (*Mocklogger)(nil)
var _ demandClient = (*MockdemandClient)(nil)
var _ httpClient = (*MockhttpClient)(nil)
var _ kafkaProducer = (*MockkafkaProducer)(nil)
var _ cacheClient = (*MockcacheClient)(nil)
var _ udws3Client = (*Mockudws3Client)(nil)
