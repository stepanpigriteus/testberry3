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
}

func ConfigMy() *Config {

	cfg := config.New()
	cfg.EnableEnv("")
	err := cfg.LoadEnvFiles(".env")
	if err != nil {
		log.Fatal("can't load .env:", err)
	}
	var cfgData Config

	cfgData.AppPort = cfg.GetString("APP_PORT")
	cfgData.DBHost = cfg.GetString("POSTGRES_HOST")
	cfgData.DBUser = cfg.GetString("POSTGRES_USER")
	cfgData.DBPass = cfg.GetString("POSTGRES_PASSWORD")
	cfgData.DBName = cfg.GetString("POSTGRES_DB")
	cfgData.DBPort = cfg.GetString("POSTGRES_PORT")
	cfgData.DBSSLMode = cfg.GetString("DB_SSLMODE")

	return &cfgData
}
