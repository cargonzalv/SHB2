package headerbidder

import (
	"testing"

	"github.com/adgear/go-commons/pkg/httpclient"
	"github.com/adgear/go-commons/pkg/kafka/producer"
	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/config"
	"github.com/adgear/sps-header-bidder/pkg/demand"
	"github.com/adgear/sps-header-bidder/pkg/kafkaclient"
	"github.com/adgear/sps-header-bidder/pkg/privacy"
	"github.com/adgear/sps-header-bidder/pkg/tifascache"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	_port              = "1234"
	_publicaUrl string = "http://test_publica_url"
)

func TestCreateServer(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := log.GetLogger()
	mockKafkaProducer := producer.NewMockService(mockCtrl)
	mockhttpClient := httpclient.NewMockService(mockCtrl)
	demandClient := demand.New(logger, _publicaUrl, mockhttpClient)
	kafkaClient := kafkaclient.New(logger, mockKafkaProducer, cfg)

	cacheClient := tifascache.NewService(logger, cfg)
	privacyClient := privacy.NewService(logger, cacheClient)

	fwNewHandler := FwNewHandler(logger, kafkaClient, demandClient, privacyClient)
	ssNewHandler := SsNewHandler(logger, kafkaClient, demandClient, privacyClient)
	handlers := &Handlers{
		springServeHandler: ssNewHandler,
		freeWheelHandler:   fwNewHandler,
	}

	headerBidder := CreateServer(logger, kafkaClient, demandClient, privacyClient, _port)
	assert.Equal(t, headerBidder.port, _port)
	assert.Equal(t, headerBidder.handlers, handlers)
	assert.Equal(t, headerBidder.logger, logger)
	assert.NotNil(t, headerBidder, "no err expected in CreateServer for valid values")
}

func TestRoutes(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := log.GetLogger()
	mockKafkaProducer := producer.NewMockService(mockCtrl)
	mockhttpClient := httpclient.NewMockService(mockCtrl)
	demandClient := demand.New(logger, _publicaUrl, mockhttpClient)
	kafkaClient := kafkaclient.New(logger, mockKafkaProducer, cfg)

	cacheClient := tifascache.NewService(logger, cfg)
	privacyClient := privacy.NewService(logger, cacheClient)

	fwNewHandler := FwNewHandler(logger, kafkaClient, demandClient, privacyClient)
	ssNewHandler := SsNewHandler(logger, kafkaClient, demandClient, privacyClient)
	handlers := &Handlers{
		springServeHandler: ssNewHandler,
		freeWheelHandler:   fwNewHandler,
	}
	route := routes(handlers)
	postBindings := route.PostBindings
	getBindings := route.GetBindings
	staticBindings := route.StaticFilesBindings

	assert.Equal(t, len(postBindings), 2)
	assert.Equal(t, len(getBindings), 2)
	assert.Equal(t, len(staticBindings), 1)

	if _, ok := postBindings["/hb"]; !ok {
		t.Error("Route missing for POST /hb")
	}

	if _, ok := getBindings["/health/liveness"]; !ok {
		t.Error("Route missing for GET /health/liveness")
	}

	if _, ok := getBindings["/health/readiness"]; !ok {
		t.Error("Route missing for GET /health/readiness")
	}

	if _, ok := staticBindings["/docs/{filepath:*}"]; !ok {
		t.Error("Route missing for static filepath /docs/{filepath:*}")
	}
}
