package pkg

import (
	"log"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	AppPort string

	KafkaBrokers string
	KafkaTopic   string
	KafkaGroupID string
	KafkaPort    string

	ZookeeperHost string
	ZookeeperPort string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
}

func ConfigMy() *Config {
	var cfgData Config
	cfg := config.New()

	err := cfg.Load(".env", "", "")
	if err != nil {
		log.Fatal("can't load .env:", err)
	}

	cfgData.AppPort = cfg.GetString("APP_PORT")

	cfgData.KafkaBrokers = cfg.GetString("KAFKA_BROKERS")
	cfgData.KafkaTopic = cfg.GetString("KAFKA_TOPIC")
	cfgData.KafkaGroupID = cfg.GetString("KAFKA_GROUP_ID")
	cfgData.KafkaPort= cfg.GetString("KAFKA_PORT_1")

	cfgData.ZookeeperHost = cfg.GetString("ZOOKEEPER_HOST")
	cfgData.ZookeeperPort = cfg.GetString("ZOOKEEPER_CLIENT_PORT")

	cfgData.MinioEndpoint = cfg.GetString("MINIO_ENDPOINT")
	cfgData.MinioAccessKey = cfg.GetString("MINIO_ROOT_USER")
	cfgData.MinioSecretKey = cfg.GetString("MINIO_ROOT_PASSWORD")
	cfgData.MinioBucket = cfg.GetString("MINIO_BUCKET")

	return &cfgData
}
