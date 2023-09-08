package main

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/adgear/go-commons/pkg/metric"
	"github.com/procyon-projects/chrono"

	"github.com/adgear/go-commons/pkg/httpclient"
	"github.com/adgear/go-commons/pkg/kafka/producer"
	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/go-commons/pkg/offlineprocessor"
	"github.com/adgear/sps-header-bidder/config"
	"github.com/adgear/sps-header-bidder/pkg/demand"
	"github.com/adgear/sps-header-bidder/pkg/headerbidder"
	"github.com/adgear/sps-header-bidder/pkg/importer"
	"github.com/adgear/sps-header-bidder/pkg/kafkaclient"
	"github.com/adgear/sps-header-bidder/pkg/privacy"
	"github.com/adgear/sps-header-bidder/pkg/tifascache"
	"github.com/adgear/sps-header-bidder/pkg/udws3"
	"github.com/oklog/run"
)

// main function will start the sps-header-bidder server and initiate
// logger, schemaregistry, config, offlinewriter, kafkaparams, metrics
// kafkaproducer and start the server
func main() {
	logger := log.GetLogger()
	metric.LoadMetricService("", "")
	// Configuration
	cfg, err := config.NewConfig("config/config.yml")
	if err != nil {
		panic(err)
	}
	logger.SetLogParams("json", cfg.Environment, cfg.Log.Level, cfg.Name, "")
	logger.Debug("Configuration loaded successfully", log.Metadata{"config": cfg})

	ensureNotForceToOffline(cfg)

	writerParams := offlineprocessor.WriterParams{
		RootPath: cfg.Kafka.OfflineStoragePath,
	}

	offlineWriter := offlineprocessor.NewWriter(logger, &writerParams)

	kafkaParams := &producer.ProducerParams{
		SchemaServerUri:       cfg.Schema.Server,
		AsyncTimeout:          3, // seconds
		LogLevel:              "info",
		EnablePool:            false,
		Brokers:               cfg.Kafka.Brokers,
		CaCert:                cfg.Ssl.CaCert,
		PrivPem:               cfg.Ssl.PrivPem,
		PubPem:                cfg.Ssl.PubPem,
		OfflineStorageEnabled: cfg.Kafka.EnableOfflineStorage,
		ForceOfflineStorage:   cfg.Kafka.ForceOfflineStorage,
		OfflineWriter:         offlineWriter,
	}

	kafkaProducer, err := producer.NewService(logger, kafkaParams)
	if err != nil {
		log.Error("Getting error while initializing kafka client", log.Metadata{"error": err, "config": cfg})
		return
	}
	cacheClient := tifascache.NewService(logger, cfg)
	udws3Client := udws3.NewService(logger, cfg)

	importerClient := importer.NewService(logger, cfg, udws3Client, cacheClient)
	taskScheduler := chrono.NewDefaultTaskScheduler()
	_, _ = taskScheduler.ScheduleAtFixedRate(func(ctx context.Context) {
		importerClient.LoadTifas()
	}, 1*time.Hour)

	privacyClient := privacy.NewService(logger, cacheClient)
	kafkaClient := kafkaclient.New(logger, kafkaProducer, cfg)
	httpClient := httpclient.NewService(logger)
	demandClient := demand.New(logger, demandURL(cfg), httpClient)
	server := headerbidder.CreateServer(logger, kafkaClient, demandClient, privacyClient, cfg.Http.Port)

	var g run.Group
	{
		ctx, cancel := context.WithCancel(context.Background())
		//start spring serve service
		g.Add(func() error {
			return server.Run(ctx)
		}, func(err error) {
			cancel()
			if err != nil {
				logger.Error("Failed to listen and serve", log.Metadata{"error": err})
			}
		})
	}
	err = g.Run()
	if err != nil {
		logger.Error("failed to run", log.Metadata{"error": err})
	}
}

// demandURL function takes Config as arguments and returns demand url
func demandURL(cfg *config.Config) string {
	baseURL := cfg.Publica.Endpoint
	v := url.Values{}
	v.Set("app_bundle", cfg.Publica.AppBundle)
	v.Set("format", cfg.Publica.Format)
	url := baseURL + "?" + v.Encode()
	return url
}

// ensureNotForceToOffline checking ForceOfflineStorage flag to ensure
// not to force offline storge
func ensureNotForceToOffline(cfg *config.Config) {
	if !cfg.Kafka.ForceOfflineStorage {
		return
	}

	if cfg.App.Environment == "prod" {
		err := errors.New("ForceOfflineStorage setting is for testing environment only, and must be disabled")
		panic(err)
	}
}
