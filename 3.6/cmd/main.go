package main

import (
	"context"
	"fmt"
	"log"
	"threeSixth/internal/db"
	"threeSixth/internal/httpsh"
	"threeSixth/internal/service"
	"threeSixth/pkg"

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
	zlog.Logger.Info().Msg("[3.1/5] Init Service")
	serv := service.NewService(ctx, zlog.Logger, *dataBase)

	zlog.Logger.Info().Msg("[3.2/5] Init Handlers")
	handlers := httpsh.NewHandlers(ctx, serv, zlog.Logger)
	zlog.Logger.Info().Msg("[3.3/4] Start Server")
	server := httpsh.NewServer(configs.AppPort, zlog.Logger, serv, handlers, dataBase)

	zlog.Logger.Info().Msg("[4/4] start server!")
	if err := server.RunServer(ctx); err != nil {
		log.Fatalf("server error: %v", err)
	}
	zlog.Logger.Info().Msg(" All components works!")

}
