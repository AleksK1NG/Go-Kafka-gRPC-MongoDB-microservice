package main

import (
	"context"
	"log"

	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/products-microservice/config"
	"github.com/AleksK1NG/products-microservice/internal/server"
	"github.com/AleksK1NG/products-microservice/pkg/jaeger"
	"github.com/AleksK1NG/products-microservice/pkg/kafka"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
	"github.com/AleksK1NG/products-microservice/pkg/mongodb"
)

// @title Products microservice
// @version 1.0
// @description Products REST API
// @contact.name Alexander Bryksin
// @contact.url https://github.com/AleksK1NG
// @contact.email alexander.bryksin@yandex.ru
// @BasePath /api/v1
func main() {
	log.Println("Starting products service")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting user server")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, DevelopmentMode: %s",
		cfg.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Development,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.AppVersion)

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, cfg)
	if err != nil {
		appLogger.Fatal("cannot connect mongodb", err)
	}
	defer func() {
		if err := mongoDBConn.Disconnect(ctx); err != nil {
			appLogger.Fatal("mongoDBConn.Disconnect", err)
		}
	}()
	appLogger.Infof("MongoDB connected: %v", mongoDBConn.NumberSessionsInProgress())

	conn, err := kafka.NewKafkaConn(cfg)
	if err != nil {
		appLogger.Fatal("NewKafkaConn", err)
	}
	defer conn.Close()
	brokers, err := conn.Brokers()
	if err != nil {
		appLogger.Fatal("conn.Brokers", err)
	}
	appLogger.Infof("Kafka connected: %v", brokers)

	s := server.NewServer(appLogger, cfg, tracer, mongoDBConn)
	appLogger.Fatal(s.Run())
}
