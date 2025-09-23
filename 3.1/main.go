package main

import (
	"fmt"

	"treeOne/domain"
	"treeOne/http"
	"treeOne/pkg"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.Init()
	zlog.Logger.Info().Msg("[1/6] Reading configuration")
	configs := pkg.ConfigMy()
	zlog.Logger.Info().Msg("[2/6] Init Postgress")
	masterDSN := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		configs.DBUser, configs.DBPass, configs.DBHost, configs.Port, configs.DBName, configs.DBSSLMode,
	)
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	slaveDSNs := []string{}
	db, err := dbpg.New(masterDSN, slaveDSNs, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("init database error %s", err)
	}
	zlog.Logger.Info().Msg("[3/6] Init Redis")
	redisConnStr := configs.Redis_host + ":" + configs.Redis_port
	client := redis.New(redisConnStr, configs.Redis_pass, configs.Redis_db)
	zlog.Logger.Info().Msg("[4/6] Init RabbitMQ")
	var handlers domain.EventHandler
	zlog.Logger.Info().Msg("[5/6] Starting Server")
	server := http.NewServer(configs.Port, zlog.Logger, db, handlers, client)

	err = server.RunServer()
	if err != nil {
		zlog.Logger.Error().Msgf("Ошибка запуска сервера: %s", err)
	}
	fmt.Println(db, client)
}
