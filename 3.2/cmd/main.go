package main

import (
	"context"
	"fmt"
	"treeTwo/internal/httpsh"
	"treeTwo/internal/httpsh/handlers"
	"treeTwo/internal/service"
	"treeTwo/internal/storage"
	"treeTwo/pkg"

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
		configs.DBUser, configs.DBPass, configs.DBHost, configs.DBPort, configs.DBName, configs.DBSSLMode,
	)
	slaveDSNs := []string{}
	storage := storage.NewStorage(ctx, masterDSN, slaveDSNs, zlog.Logger)
	zlog.Logger.Info().Msg("[3/6] Init Redis")
	redisConnStr := configs.Redis_host + ":" + configs.Redis_port
	client := redis.New(redisConnStr, configs.Redis_pass, configs.Redis_db)
	zlog.Logger.Info().Msg("[4.1/6] Init Service")
	serv := service.NewService(ctx, storage, zlog.Logger, *client)
	zlog.Logger.Info().Msg("[4.2/6] Init Handlers")
	handlers := handlers.NewHandlers(ctx, serv, zlog.Logger)
	zlog.Logger.Info().Msg("[4.3/6] Start Server")
	server := httpsh.NewServer(configs.Port, zlog.Logger, serv, storage, handlers, client)
	server.RunServer()

}
