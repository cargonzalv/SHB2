// Package kafkaclient - .
package kafkaclient

import (
	"github.com/adgear/go-commons/pkg/kafka/producer"
	"github.com/adgear/go-commons/pkg/log"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/adgear/sps-header-bidder/config"
)

// EventInfo parameters.
type EventInfo struct {
	// producer kafka topic
	topic string
	// schema key
	schemaKey string
	// schema file path
	schemaFile string
}

type kafkaClient struct {
	// logger service
	logger log.Service
	// kafka producer service
	kafka producer.Service
	// ad Event info struct
	adEvent *EventInfo
	// inventtory Event info struct
	inventoryEvent *EventInfo
}

// This statement forcing the module to
// implements the kafkaClient `Service` interface.
var _ Service = (*kafkaClient)(nil)

// New is a constructor function which get logger implementation and
// kafkaClient params arguments and return implementation of KafkaClient interface.
func New(l log.Service, k producer.Service, c *config.Config) *kafkaClient {
	return &kafkaClient{
		logger: l,
		kafka:  k,
		adEvent: &EventInfo{
			topic:      c.Kafka.AdEventTopic,
			schemaKey:  c.Schema.AdEventKey,
			schemaFile: c.Schema.AdEventSchemaFile,
		},
		inventoryEvent: &EventInfo{
			topic:      c.Kafka.InventoryEventTopic,
			schemaKey:  c.Schema.InventoryEventKey,
			schemaFile: c.Schema.InventoryEventSchemaFile,
		},
	}
}

// PublishAdEventAsync publishes asynchronous Ad events.
func (k *kafkaClient) PublishAdEventAsync(message []byte) {
	k.publishEventAsync(k.adEvent.topic, message, k.adEvent.schemaKey, k.adEvent.schemaFile)
}

// PublishInventoryEventAsync publishes asynchronous Inventory events.
func (k *kafkaClient) PublishInventoryEventAsync(message []byte) {
	k.publishEventAsync(k.inventoryEvent.topic, message, k.inventoryEvent.schemaKey, k.inventoryEvent.schemaFile)
}

// publishEventAsync publishes events asynchronously to kafka
func (k *kafkaClient) publishEventAsync(
	topic string,
	message []byte,
	schemaKey string,
	schemaFile string) {
	metadata := log.Metadata{
		"topic":      topic,
		"message":    string(message),
		"schemaKey":  schemaKey,
		"schemaFile": schemaFile,
	}
	k.logger.Debug("Publishing event to kafka", metadata)
	promise := func(record *kgo.Record, err error) {
		if err != nil {
			k.logger.Error("Get an error while producing kafka event", metadata)
			return
		}
		k.logger.Debug("Completed publish kafka event with success", metadata)
	}
	k.kafka.PublishKafkaEvent(topic, message, schemaKey, schemaFile, promise)
}
