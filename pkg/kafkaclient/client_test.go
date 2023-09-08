package kafkaclient

import (
	"sync"
	"testing"

	"github.com/adgear/go-commons/pkg/kafka/producer"
	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestKafkaClientNew(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := log.GetLogger()
	mockKafkaProducer := producer.NewMockService(mockCtrl)
	kafkaClient := New(logger, mockKafkaProducer, cfg)
	assert.NotNil(t, kafkaClient.adEvent, "adEvent is configured")
	assert.NotNil(t, kafkaClient.inventoryEvent, "inventoryEvent is configured")
	assert.NotNil(t, kafkaClient, "no err expected in kafkaClient for valid values")
}

func TestKafkaClientPublishAdEventAsync(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := log.GetLogger()
	mockKafkaProducer := producer.NewMockService(mockCtrl)
	kafkaClient := New(logger, mockKafkaProducer, cfg)

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(cfg.Kafka.AdEventTopic, gomock.Any(), cfg.Schema.AdEventKey, cfg.Schema.AdEventSchemaFile, gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		})
	kafkaClient.PublishAdEventAsync([]byte("{}"))
}

func TestKafkaClientPublishInventoryEventAsync(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := log.GetLogger()
	mockKafkaProducer := producer.NewMockService(mockCtrl)
	kafkaClient := New(logger, mockKafkaProducer, cfg)

	mockKafkaProducer.EXPECT().
		PublishKafkaEvent(cfg.Kafka.InventoryEventTopic, gomock.Any(), cfg.Schema.InventoryEventKey, cfg.Schema.InventoryEventSchemaFile, gomock.Any()).
		Do(func(arg0, arg1, arg2, arg3, arg4 interface{}) {
			defer wg.Done()
		})
	kafkaClient.PublishInventoryEventAsync([]byte("{}"))
}
