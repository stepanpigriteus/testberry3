package pkg

import (
	"log"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	AppPort string

	DBUser    string
	DBPass    string
	DBHost    string
	DBPort    string
	DBName    string
	DBSSLMode string

	KafkaBrokers string
	KafkaTopic   string
	KafkaGroupID string
	KafkaPort    string

	ZookeeperHost string
	ZookeeperPort string
}

func ConfigMy() *Config {

	cfg := config.New()
	cfg.EnableEnv("")
	err := cfg.LoadEnvFiles(".env")
	if err != nil {
		log.Fatal("can't load .env:", err)
	}
	var cfgData Config
	if cfgData.KafkaBrokers == "" {
		cfgData.KafkaBrokers = "kafka:" + cfgData.KafkaPort
	}
	cfgData.AppPort = cfg.GetString("APP_PORT")
	cfgData.DBHost = cfg.GetString("POSTGRES_HOST")
	cfgData.DBUser = cfg.GetString("POSTGRES_USER")
	cfgData.DBPass = cfg.GetString("POSTGRES_PASSWORD")
	cfgData.DBName = cfg.GetString("POSTGRES_DB")
	cfgData.DBPort = cfg.GetString("POSTGRES_PORT")
	cfgData.DBSSLMode = cfg.GetString("DB_SSLMODE")

	if cfgData.KafkaBrokers == "" {
		cfgData.KafkaBrokers = "kafka:" + cfgData.KafkaPort
	} else {
		cfgData.KafkaBrokers = cfg.GetString("KAFKA_BROKERS")
	}
	cfgData.KafkaTopic = cfg.GetString("KAFKA_TOPIC")
	cfgData.KafkaGroupID = cfg.GetString("KAFKA_GROUP_ID")
	cfgData.KafkaPort = cfg.GetString("KAFKA_PORT_1")

	cfgData.ZookeeperHost = cfg.GetString("ZOOKEEPER_HOST")
	cfgData.ZookeeperPort = cfg.GetString("ZOOKEEPER_CLIENT_PORT")

	return &cfgData
}
