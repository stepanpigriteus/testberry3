package main

import (
	"fmt"
	"threeFour/pkg"

	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	// ctx := context.Background()
	zlog.Init()
	zlog.Logger.Info().Msg("[1/6] Reading configuration")
	configs := pkg.ConfigMy()
	fmt.Println(configs)
	// zlog.Logger.Info().Msg("[2/6] Init Postgress")

	zlog.Logger.Info().Msg("[3/6] Init Kafka")
	_ = kafka.NewProducer([]string{"localhost:" + configs.KafkaPort}, "topic")
	// zlog.Logger.Info().Msg("[4.1/5] Init Service")
	// serv := service.NewService(ctx, storage, zlog.Logger)
	// zlog.Logger.Info().Msg("[4.2/5] Init Handlers")
	// handlers := httpsh.NewHandlers(ctx, serv, zlog.Logger)
	// zlog.Logger.Info().Msg("[4.3/5] Start Server")
	// server := httpsh.NewServer(configs.Port, zlog.Logger, serv, storage, handlers)
	// server.RunServer()
	zlog.Logger.Info().Msg("[5/5] All components works!")
}
