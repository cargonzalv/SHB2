package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App        `yaml:"app"`
		Http       `yaml:"http"`
		Connection `yaml:"connection"`
		Log        `yaml:"logger"`
		Kafka      `yaml:"kafka"`
		Schema     `yaml:"schema_registry"`
		Ssl        `yaml:"ssl"`
		Publica    `yaml:"publica"`
		S3         `yaml:"s3"`
		Cache      `yaml:"cache"`
	}

	// App -.
	App struct {
		Name        string `env-required:"true" yaml:"name"        env:"HB_APP_NAME"`
		Version     string `env-required:"true" yaml:"version"     env:"HB_APP_VERSION"`
		Environment string `env-required:"true" yaml:"environment" env:"ENVIRONMENT"`
	}

	// Http -.
	Http struct {
		Host string `env-required:"true" yaml:"host" env:"SERVER_HOST"`
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}

	// Connection -.
	Connection struct {
		Retry   int `env-required:"true" yaml:"retry"   env:"CONNECTION_RETRY"`
		Timeout int `env-required:"true" yaml:"timeout" env:"CONNECTION_TIMEOUT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	// Kafka -.
	Kafka struct {
		EnableOfflineStorage bool   `env-required:"true" yaml:"enable_offline_storage" env:"KAFKA_ENABLE_OFFLINE_STORAGE"`
		ForceOfflineStorage  bool   `env-required:"false" yaml:"force_offline_storage" env:"KAFKA_FORCE_OFFLINE_STORAGE"`
		OfflineStoragePath   string `env-required:"false" yaml:"offline_storage_path"   env:"KAFKA_OFFLINE_STORAGE_PATH"`
		InventoryEventTopic  string `env-required:"true" yaml:"inventory_event_topic"  env:"KAFKA_INVENTORY_EVENT_TOPIC"`
		AdEventTopic         string `env-required:"true" yaml:"ad_event_topic"         env:"KAFKA_AD_EVENT_TOPIC"`
		Brokers              string `env-required:"true" yaml:"brokers"                env:"KAFKA_BROKERS"`
		PoolSize             int    `env-required:"true" yaml:"pool_size"              env:"KAFKA_POOL_SIZE"`
	}

	// Schema -.
	Schema struct {
		Server                   string `env-required:"true" yaml:"server"                      env:"SCHEMA_SERVER"`
		AdEventKey               string `env-required:"true" yaml:"ad_event_key"                env:"SCHEMA_AD_EVENT_KEY"`
		InventoryEventKey        string `env-required:"true" yaml:"inventory_event_key"         env:"SCHEMA_INVENTORY_EVENT_KEY"`
		AdEventSchemaFile        string `env-required:"true" yaml:"ad_event_schema_file"        env:"SCHEMA_AD_EVENT_SCHEMA_FILE"`
		InventoryEventSchemaFile string `env-required:"true" yaml:"inventory_event_schema_file" env:"SCHEMA_INVENTORY_EVENT_SCHEMA_FILE"`
	}

	// Ssl -.
	Ssl struct {
		CaCert  string `env-required:"true" yaml:"ca_cert"  env:"CA_CERT"`
		PrivPem string `env-required:"true" yaml:"priv_pem" env:"PRIV_PEM"`
		PubPem  string `env-required:"true" yaml:"pub_pem"  env:"PUB_PEM"`
	}

	// Publica -.
	Publica struct {
		Endpoint  string `env-required:"true" yaml:"endpoint"   env:"PUBLICA_ENDPOINT"`
		AppBundle string `env-required:"true" yaml:"app_bundle" env:"PUBLICA_APP_BUNDLE"`
		Format    string `env-required:"true" yaml:"format"     env:"PUBLICA_FORMAT"`
	}
	// S3 -.
	S3 struct {
		Region               string `env-required:"true" yaml:"region" env:"S3_REGION"`
		Bucket               string `env-required:"true" yaml:"bucket" env:"S3_BUCKET"`
		AssumeRole           string `env-required:"true" yaml:"assume_role" env:"ASSUME_ROLE"`
		OptoutPrefix         string `env-required:"true" yaml:"optout_prefix" env:"OPTOUT_PREFIX"`
		OptputSuccessFileKey string `env-required:"true" yaml:"optout_success_file_key" env:"OPTOUT_SUCCESS_FILE_KEY"`
	}

	// Cache -.
	Cache struct {
		KeyTTL           string `env-required:"true" yaml:"key_ttl" env:"TIFA_KEY_TTL"`
		LoadTimestampKey string `env-required:"true" yaml:"load_timestamp_key" env:"LAST_LOAD_TIMESTAMP_KEY"`
	}
)

// NewConfig returns app config.
func NewConfig(path ...string) (*Config, error) {
	var configPath string
	if len(path) > 0 && len(path[0]) > 0 {
		configPath = path[0]
	} else {
		configPath = "config.yml"
	}

	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)

	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	loadEnv(cfg)

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadEnv(cfg *Config) {
	env := cfg.Environment
	_ = godotenv.Load(".env/.env." + env + ".local")
	_ = godotenv.Load(".env/.env." + env)
	_ = godotenv.Load(".env/.env")
}
