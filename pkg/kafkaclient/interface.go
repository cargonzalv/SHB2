// Kafka Client package.
package kafkaclient

// `interface.go` module contains declaration of interfaces,
// provided by `kafkaclient` package
type Service interface {
	PublishAdEventAsync(message []byte)
	PublishInventoryEventAsync(message []byte)
}
