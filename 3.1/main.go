package main

import (
	"context"
	"fmt"
	"treeOne/http"
	"treeOne/pkg"
	"treeOne/service"
	"treeOne/storage"

	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	ctx := context.Background()
	zlog.Init()
	zlog.Logger.Info().Msg("[1/6] Reading configuration")
	configs := pkg.ConfigMy()
	fmt.Println(configs)
	zlog.Logger.Info().Msg("[2/6] Init Postgress")
	masterDSN := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		configs.DBUser, configs.DBPass, configs.DBHost, configs.Port, configs.DBName, configs.DBSSLMode,
	)

	slaveDSNs := []string{}
	storage := storage.NewStorage(ctx, masterDSN, slaveDSNs, zlog.Logger)

	zlog.Logger.Info().Msg("[3/6] Init Redis")
	redisConnStr := configs.Redis_host + ":" + configs.Redis_port
	client := redis.New(redisConnStr, configs.Redis_pass, configs.Redis_db)
	zlog.Logger.Info().Msg("[3.1/6] Init Service")
	service := service.NewService(storage, zlog.Logger, *client)
	zlog.Logger.Info().Msg("[3.2/6] Init Handlers")
	handlers := http.NewHandleNotify(zlog.Logger, service, *client)
	zlog.Logger.Info().Msg("[4/6] Init RabbitMQ")
	zlog.Logger.Info().Msg("[5/6] Starting Server")

	server := http.NewServer(configs.Port, zlog.Logger, service, storage, handlers, client)

	err := server.RunServer()
	if err != nil {
		zlog.Logger.Error().Msgf("Ошибка запуска сервера: %s", err)
	}
	fmt.Println(client)
}
