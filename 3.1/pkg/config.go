package pkg

import "github.com/wb-go/wbf/config"

type Config struct {
	Port      string
	DBUser    string
	DBPass    string
	DBHost    string
	DBPort    string
	DBName    string
	DBSSLMode string

	Redis_host    string
	Redis_port    string
	Redis_pass    string
	Redis_db      int
	Rabb_user     string
	Rabb_pass     string
	Rabb_exchange string
}

func ConfigMy() *Config {
	var configs Config
	cfg := config.New()
	cfg.Load(".env")

	configs.Port = cfg.GetString("PORT")
	configs.DBHost = cfg.GetString("DB_HOST")
	configs.DBUser = cfg.GetString("DB_USER")
	configs.DBPass = cfg.GetString("DB_PASSWORD")
	configs.DBName = cfg.GetString("DB_NAME")
	configs.DBPort = cfg.GetString("DB_PORT")
	configs.DBSSLMode = cfg.GetString("DB_SSLMODE")
	configs.Redis_host = cfg.GetString("REDIS_HOST")
	configs.Redis_port = cfg.GetString("REDIS_PORT")
	configs.Redis_pass = cfg.GetString("REDIS_PASS")
	configs.Redis_db = cfg.GetInt("REDIS_DB")

	configs.Rabb_user = cfg.GetString("RABBITMQ_USER")
	configs.Rabb_pass = cfg.GetString("RABBITMQ_PASS")
	configs.Rabb_exchange = cfg.GetString("RABBITMQ_EXCHANGE")

	return &configs
}
