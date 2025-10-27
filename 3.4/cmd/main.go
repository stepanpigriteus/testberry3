package main

import (
	"context"
	"fmt"
	"log"
	"threeFour/internal/db"
	"threeFour/internal/httpsh"
	"threeFour/internal/service"
	"threeFour/internal/storage"
	"threeFour/pkg"

	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	ctx := context.Background()
	zlog.Init()
	zlog.Logger.Info().Msg("[1/6] Reading configuration")
	configs := pkg.ConfigMy()
	zlog.Logger.Info().Msg("[2/6] Init Postgress")
	masterDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		configs.DBHost, configs.DBPort, configs.DBUser, configs.DBPass, configs.DBName, configs.DBSSLMode,
	)
	slaveDSNs := []string{}
	dataBase := db.NewDb(ctx, masterDSN, slaveDSNs, zlog.Logger)

	zlog.Logger.Info().Msg("[2/6] Init Minio")
	storage, err := storage.NewMinioStorage(configs.MinioEndpoint, configs.MinioAccessKey, configs.MinioSecretKey, false)
	if err != nil {
		log.Fatal(err)
	}
	zlog.Logger.Info().Msg("[3/6] Init Kafka")
	err = pkg.CreateTopic([]string{configs.KafkaBrokers}, configs.KafkaTopic)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to create Kafka topic")
	}
	producer := kafka.NewProducer([]string{configs.KafkaBrokers}, configs.KafkaTopic)

	consumer := kafka.NewConsumer([]string{configs.KafkaBrokers}, configs.KafkaTopic, configs.KafkaGroupID)

	zlog.Logger.Info().Msg("[4.1/5] Init Service")
	serv := service.NewService(ctx, producer, consumer, zlog.Logger, storage, *dataBase)
	go serv.Convee(ctx)
	zlog.Logger.Info().Msg("[4.2/5] Init Handlers")
	handlers := httpsh.NewHandlers(ctx, serv, zlog.Logger)
	zlog.Logger.Info().Msg("[4.3/5] Start Server")
	server := httpsh.NewServer(configs.AppPort, zlog.Logger, serv, handlers, dataBase, storage)
	server.RunServer()
	zlog.Logger.Info().Msg("[5/5] All components works!")

}
