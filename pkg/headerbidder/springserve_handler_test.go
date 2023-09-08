package headerbidder

import (
	"io"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/adgear/sps-header-bidder/config"
	"github.com/adgear/sps-header-bidder/pkg/headerbidder/apitest"
	"github.com/adgear/sps-header-bidder/pkg/kafkaclient"
	"github.com/adgear/sps-header-bidder/pkg/mockservices"
	"github.com/adgear/sps-header-bidder/pkg/privacy"
	"github.com/adgear/sps-header-bidder/pkg/tifascache"
)

// Happy path
func TestSpringServeRequestHandlerHappyPath(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	responseNoContent := []byte{}
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mockservices.NewMocklogger(mockCtrl)
	mockKafkaProducer := mockservices.NewMockkafkaProducer(mockCtrl)
	mockDemand := mockservices.NewMockdemandClient(mockCtrl)

	cacheClient := tifascache.NewService(mockLogger, cfg)
	privacyClient := privacy.NewService(mockLogger, cacheClient)

	mockKafka := kafkaclient.New(mockLogger, mockKafkaProducer, cfg)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(cfg.Kafka.AdEventTopic, gomock.Any(), cfg.Schema.AdEventKey, cfg.Schema.AdEventSchemaFile, gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		})

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(
			cfg.Kafka.InventoryEventTopic,
			gomock.Any(),
			cfg.Schema.InventoryEventKey,
			cfg.Schema.InventoryEventSchemaFile,
			gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		})

	mockDemand.EXPECT().BidOrtbReq(gomock.Any()).Return(204, responseNoContent)

	req, _ := apitest.Request(apitest.Samples().SpringserveApiRequests.Valid, apitest.Samples().Url)
	handler := SsNewHandler(mockLogger, mockKafka, mockDemand, privacyClient)
	listener := apitest.InitServer(handler.Handler)
	resp, err := apitest.ProcessRequest(listener, req, false)
	respBody, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 204)
	assert.Equal(t, respBody, responseNoContent)
	assert.NoError(t, err)
	wg.Wait()
}

// Invalid demand response for springserve Req
func TestSpringServeRequestHandlerInvalidDemandResp(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mockservices.NewMocklogger(mockCtrl)
	mockKafkaProducer := mockservices.NewMockkafkaProducer(mockCtrl)
	mockDemand := mockservices.NewMockdemandClient(mockCtrl)

	cacheClient := tifascache.NewService(mockLogger, cfg)
	privacyClient := privacy.NewService(mockLogger, cacheClient)

	mockLogger.
		EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()
	mockLogger.
		EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	mockKafkaProducer.
		EXPECT().
		PublishKafkaEvent(
			cfg.Kafka.AdEventTopic,
			gomock.Any(),
			cfg.Schema.AdEventKey,
			cfg.Schema.AdEventSchemaFile,
			gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		})

	mockKafkaProducer.
		EXPECT().
		PublishKafkaEvent(
			cfg.Kafka.InventoryEventTopic,
			gomock.Any(),
			cfg.Schema.InventoryEventKey,
			cfg.Schema.InventoryEventSchemaFile,
			gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		})

	mockDemand.EXPECT().BidOrtbReq(gomock.Any()).Return(500, []byte("invalid response"))

	mockKafka := kafkaclient.New(mockLogger, mockKafkaProducer, cfg)

	req, _ := apitest.Request(apitest.Samples().SpringserveApiRequests.Valid, apitest.Samples().Url)
	handler := SsNewHandler(mockLogger, mockKafka, mockDemand, privacyClient)
	listener := apitest.InitServer(handler.Handler)
	resp, err := apitest.ProcessRequest(listener, req, false)

	assert.Equal(t, resp.StatusCode, 500)
	assert.NoError(t, err)
	wg.Wait()
}

// Invalid request body for Springserve request
func TestSpringServeRequestHandlerInvalidRequestBody(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mockservices.NewMocklogger(mockCtrl)
	mockKafkaProducer := mockservices.NewMockkafkaProducer(mockCtrl)
	mockDemand := mockservices.NewMockdemandClient(mockCtrl)

	cacheClient := tifascache.NewService(mockLogger, cfg)
	privacyClient := privacy.NewService(mockLogger, cacheClient)

	mockLogger.
		EXPECT().
		Debug(gomock.Any(), gomock.Any())
	mockLogger.
		EXPECT().
		Error(gomock.Any(), gomock.Any())

	mockKafka := kafkaclient.New(mockLogger, mockKafkaProducer, cfg)

	req, _ := apitest.Request(apitest.Samples().SpringserveApiRequests.InvalidJson, apitest.Samples().Url)
	handler := SsNewHandler(mockLogger, mockKafka, mockDemand, privacyClient)
	listener := apitest.InitServer(handler.Handler)
	resp, err := apitest.ProcessRequest(listener, req, false)
	assert.Equal(t, resp.StatusCode, 400)
	assert.NoError(t, err)
}
