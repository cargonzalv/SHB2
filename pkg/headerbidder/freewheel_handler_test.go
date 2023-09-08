package headerbidder

import (
	"io"
	"sync"
	"testing"

	"github.com/adgear/sps-header-bidder/config"
	"github.com/adgear/sps-header-bidder/pkg/demand"
	"github.com/adgear/sps-header-bidder/pkg/headerbidder/apitest"
	"github.com/adgear/sps-header-bidder/pkg/kafkaclient"
	"github.com/adgear/sps-header-bidder/pkg/mockservices"
	"github.com/adgear/sps-header-bidder/pkg/privacy"
	"github.com/adgear/sps-header-bidder/pkg/tifascache"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

// Happy path - Free Wheel header bidder request with valid "Deal Id"
// Expect status code 200. Response include deal id from the request
func TestFreeWheelRequestWithDealIdHandlerHappyPath(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mockservices.NewMocklogger(mockCtrl)
	mockKafkaProducer := mockservices.NewMockkafkaProducer(mockCtrl)
	mockDemand := mockservices.NewMockdemandClient(mockCtrl)
	mockKafka := kafkaclient.New(mockLogger, mockKafkaProducer, cfg)

	cacheClient := tifascache.NewService(mockLogger, cfg)
	privacyClient := privacy.NewService(mockLogger, cacheClient)

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

	requestBody := apitest.Samples().FreewheelApiRequests.ValidDealId
	demandResponseBody := apitest.Samples().DemandApiResponse.ValidWithDealId

	expectedDemandParams := demand.DemandExtParams{
		Body:     requestBody,
		SamFlag:  "pbid",
		SamHbTag: "27083",
	}

	mockDemand.EXPECT().BidOrtbReq(expectedDemandParams).Return(200, demandResponseBody)

	req, _ := apitest.Request(requestBody, apitest.Samples().Url)
	handler := FwNewHandler(mockLogger, mockKafka, mockDemand, privacyClient)
	listener := apitest.InitServer(handler.Handler)
	resp, err := apitest.ProcessRequest(listener, req, false)
	respBody, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 200)
	assert.NotEqual(t, len(respBody), 0)
	respActualDealId := gjson.GetBytes(respBody, "seatbid.0.bid.0.dealid").String()
	assert.Equal(t, respActualDealId, "ATN-000000")
	assert.NoError(t, err)
}

// Freewheel header bidder request has empty deals
// Expect status code 204 with empty body
func TestFreeWheelRequestWithoutDealsHandler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	noContent := []byte{}
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mockservices.NewMocklogger(mockCtrl)
	mockKafkaProducer := mockservices.NewMockkafkaProducer(mockCtrl)
	mockDemand := mockservices.NewMockdemandClient(mockCtrl)
	mockKafka := kafkaclient.New(mockLogger, mockKafkaProducer, cfg)

	cacheClient := tifascache.NewService(mockLogger, cfg)
	privacyClient := privacy.NewService(mockLogger, cacheClient)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(cfg.Kafka.AdEventTopic, gomock.Any(), cfg.Schema.AdEventKey, cfg.Schema.AdEventSchemaFile, gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		}).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(
			cfg.Kafka.InventoryEventTopic,
			gomock.Any(),
			cfg.Schema.InventoryEventKey,
			cfg.Schema.InventoryEventSchemaFile,
			gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		}).AnyTimes()

	requestBody := apitest.Samples().FreewheelApiRequests.NoDeals

	expectedDemandParams := demand.DemandExtParams{
		Body:     noContent,
		SamFlag:  "",
		SamHbTag: "",
	}

	mockDemand.EXPECT().BidOrtbReq(expectedDemandParams).Return(204, noContent)

	req, _ := apitest.Request(requestBody, apitest.Samples().Url)
	handler := FwNewHandler(mockLogger, mockKafka, mockDemand, privacyClient)
	listener := apitest.InitServer(handler.Handler)
	resp, err := apitest.ProcessRequest(listener, req, false)
	respBody, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 204)
	assert.Equal(t, respBody, noContent)
	assert.NoError(t, err)
}

// Freewheel header bidder request has empty deal Id
// Expect status code 204 with empty body
func TestFreeWheelRequestWithEmptyDealIdHandler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	noContent := []byte{}
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mockservices.NewMocklogger(mockCtrl)
	mockKafkaProducer := mockservices.NewMockkafkaProducer(mockCtrl)
	mockDemand := mockservices.NewMockdemandClient(mockCtrl)
	mockKafka := kafkaclient.New(mockLogger, mockKafkaProducer, cfg)

	cacheClient := tifascache.NewService(mockLogger, cfg)
	privacyClient := privacy.NewService(mockLogger, cacheClient)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(cfg.Kafka.AdEventTopic, gomock.Any(), cfg.Schema.AdEventKey, cfg.Schema.AdEventSchemaFile, gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		}).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(
			cfg.Kafka.InventoryEventTopic,
			gomock.Any(),
			cfg.Schema.InventoryEventKey,
			cfg.Schema.InventoryEventSchemaFile,
			gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		}).AnyTimes()

	requestBody := apitest.Samples().FreewheelApiRequests.DealsWithEmptyId

	expectedDemandParams := demand.DemandExtParams{
		Body:     noContent,
		SamFlag:  "",
		SamHbTag: "",
	}

	mockDemand.EXPECT().BidOrtbReq(expectedDemandParams).Return(204, noContent)

	req, _ := apitest.Request(requestBody, apitest.Samples().Url)
	handler := FwNewHandler(mockLogger, mockKafka, mockDemand, privacyClient)
	listener := apitest.InitServer(handler.Handler)
	resp, err := apitest.ProcessRequest(listener, req, false)
	respBody, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 204)
	assert.Equal(t, respBody, noContent)
	assert.NoError(t, err)
}

// Freewheel header bidder request has valid deal Id
// Demand client return the responce without dealid
// Expect status code 200. Response include deal id from the request
func TestFreeWheelRequestWhenDemandResponseWithoutDealId(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mockservices.NewMocklogger(mockCtrl)
	mockKafkaProducer := mockservices.NewMockkafkaProducer(mockCtrl)
	mockDemand := mockservices.NewMockdemandClient(mockCtrl)
	mockKafka := kafkaclient.New(mockLogger, mockKafkaProducer, cfg)

	cacheClient := tifascache.NewService(mockLogger, cfg)
	privacyClient := privacy.NewService(mockLogger, cacheClient)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(cfg.Kafka.AdEventTopic, gomock.Any(), cfg.Schema.AdEventKey, cfg.Schema.AdEventSchemaFile, gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		}).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(
			cfg.Kafka.InventoryEventTopic,
			gomock.Any(),
			cfg.Schema.InventoryEventKey,
			cfg.Schema.InventoryEventSchemaFile,
			gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		}).AnyTimes()

	requestBody := apitest.Samples().FreewheelApiRequests.ValidDealId
	demandResponseBody := apitest.Samples().DemandApiResponse.ValidWithoutDealId

	expectedDemandParams := demand.DemandExtParams{
		Body:     requestBody,
		SamFlag:  "pbid",
		SamHbTag: "27083",
	}

	mockDemand.EXPECT().BidOrtbReq(expectedDemandParams).Return(200, demandResponseBody)

	req, _ := apitest.Request(requestBody, apitest.Samples().Url)
	handler := FwNewHandler(mockLogger, mockKafka, mockDemand, privacyClient)
	listener := apitest.InitServer(handler.Handler)
	resp, err := apitest.ProcessRequest(listener, req, false)
	respBody, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 200)
	assert.NotEqual(t, len(respBody), 0)
	respActualDealId := gjson.GetBytes(respBody, "seatbid.0.bid.0.dealid").String()
	assert.Equal(t, respActualDealId, "ATN-000000")
	assert.NoError(t, err)
}

// Freewheel header bidder request has valid deal Id
// Demand client return the invlid responce
// Expect status code 200. Demand responce return as is
func TestFreeWheelRequestWhenDemandInvalidResponse(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mockservices.NewMocklogger(mockCtrl)
	mockKafkaProducer := mockservices.NewMockkafkaProducer(mockCtrl)
	mockDemand := mockservices.NewMockdemandClient(mockCtrl)
	mockKafka := kafkaclient.New(mockLogger, mockKafkaProducer, cfg)

	cacheClient := tifascache.NewService(mockLogger, cfg)
	privacyClient := privacy.NewService(mockLogger, cacheClient)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(cfg.Kafka.AdEventTopic, gomock.Any(), cfg.Schema.AdEventKey, cfg.Schema.AdEventSchemaFile, gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		}).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(
			cfg.Kafka.InventoryEventTopic,
			gomock.Any(),
			cfg.Schema.InventoryEventKey,
			cfg.Schema.InventoryEventSchemaFile,
			gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		}).AnyTimes()

	requestBody := apitest.Samples().FreewheelApiRequests.ValidDealId
	demandUnexpectedResponseBody := apitest.Samples().DemandApiResponse.Unexpected

	expectedDemandParams := demand.DemandExtParams{
		Body:     requestBody,
		SamFlag:  "pbid",
		SamHbTag: "27083",
	}

	mockDemand.EXPECT().BidOrtbReq(expectedDemandParams).Return(200, demandUnexpectedResponseBody)

	req, _ := apitest.Request(requestBody, apitest.Samples().Url)
	handler := FwNewHandler(mockLogger, mockKafka, mockDemand, privacyClient)
	listener := apitest.InitServer(handler.Handler)
	resp, err := apitest.ProcessRequest(listener, req, false)
	respBody, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, respBody, demandUnexpectedResponseBody)
	assert.NoError(t, err)
}

// Invalid demand response for Freewheel Request
func TestFreeWheelRequestHandlerInvalidDemandResp(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mockservices.NewMocklogger(mockCtrl)
	mockKafkaProducer := mockservices.NewMockkafkaProducer(mockCtrl)
	mockDemand := mockservices.NewMockdemandClient(mockCtrl)
	mockKafka := kafkaclient.New(mockLogger, mockKafkaProducer, cfg)

	cacheClient := tifascache.NewService(mockLogger, cfg)
	privacyClient := privacy.NewService(mockLogger, cacheClient)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(cfg.Kafka.AdEventTopic, gomock.Any(), cfg.Schema.AdEventKey, cfg.Schema.AdEventSchemaFile, gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		}).AnyTimes()

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(
			cfg.Kafka.InventoryEventTopic,
			gomock.Any(),
			cfg.Schema.InventoryEventKey,
			cfg.Schema.InventoryEventSchemaFile,
			gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		}).AnyTimes()

	requestBody := apitest.Samples().FreewheelApiRequests.ValidDealId
	invalidDemandResponse := apitest.Samples().DemandApiResponse.Unexpected

	expectedDemandParams := demand.DemandExtParams{
		Body:     requestBody,
		SamFlag:  "pbid",
		SamHbTag: "27083",
	}

	mockDemand.EXPECT().BidOrtbReq(expectedDemandParams).Return(500, invalidDemandResponse)

	req, _ := apitest.Request(requestBody, apitest.Samples().Url)
	handler := FwNewHandler(mockLogger, mockKafka, mockDemand, privacyClient)
	listener := apitest.InitServer(handler.Handler)
	resp, err := apitest.ProcessRequest(listener, req, false)

	assert.Equal(t, resp.StatusCode, 500)
	assert.NoError(t, err)
}

// Invalid request body for Freewheel request
func TestFreeWheelRequestHandlerInvalidRequestBody(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mockservices.NewMocklogger(mockCtrl)
	mockKafkaProducer := mockservices.NewMockkafkaProducer(mockCtrl)
	mockDemand := mockservices.NewMockdemandClient(mockCtrl)

	cacheClient := tifascache.NewService(mockLogger, cfg)
	privacyClient := privacy.NewService(mockLogger, cacheClient)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	mockKafka := kafkaclient.New(mockLogger, mockKafkaProducer, cfg)

	invalidRequestBody := apitest.Samples().FreewheelApiRequests.InvalidJson

	req, _ := apitest.Request(invalidRequestBody, apitest.Samples().Url)
	handler := FwNewHandler(mockLogger, mockKafka, mockDemand, privacyClient)
	listener := apitest.InitServer(handler.Handler)
	resp, err := apitest.ProcessRequest(listener, req, false)
	assert.Equal(t, resp.StatusCode, 400)
	assert.NoError(t, err)
}
