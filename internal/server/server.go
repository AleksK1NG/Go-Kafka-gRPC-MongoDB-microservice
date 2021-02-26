package server

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/AleksK1NG/products-microservice/config"
	"github.com/AleksK1NG/products-microservice/internal/interceptors"
	"github.com/AleksK1NG/products-microservice/internal/middlewares"
	product "github.com/AleksK1NG/products-microservice/internal/product/delivery/grpc"
	productsHttpV1 "github.com/AleksK1NG/products-microservice/internal/product/delivery/http/v1"
	"github.com/AleksK1NG/products-microservice/internal/product/delivery/kafka"
	"github.com/AleksK1NG/products-microservice/internal/product/repository"
	"github.com/AleksK1NG/products-microservice/internal/product/usecase"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
	productsService "github.com/AleksK1NG/products-microservice/proto/product"
)

const (
	certFile        = "ssl/server.crt"
	keyFile         = "ssl/server.pem"
	maxHeaderBytes  = 1 << 20
	gzipLevel       = 5
	stackSize       = 1 << 10 // 1 KB
	csrfTokenHeader = "X-CSRF-Token"
	bodyLimit       = "2M"
	kafkaGroupID    = "products_group"
)

// server
type server struct {
	log     logger.Logger
	cfg     *config.Config
	tracer  opentracing.Tracer
	mongoDB *mongo.Client
	echo    *echo.Echo
	redis   *redis.Client
}

// NewServer constructor
func NewServer(log logger.Logger, cfg *config.Config, tracer opentracing.Tracer, mongoDB *mongo.Client, redis *redis.Client) *server {
	return &server{log: log, cfg: cfg, tracer: tracer, mongoDB: mongoDB, echo: echo.New(), redis: redis}
}

// Run Start server
func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	validate := validator.New()

	productsProducer := kafka.NewProductsProducer(s.log, s.cfg)
	productsProducer.Run()
	defer productsProducer.Close()

	productMongoRepo := repository.NewProductMongoRepo(s.mongoDB)
	productRedisRepo := repository.NewProductRedisRepository(s.redis)
	productUC := usecase.NewProductUC(productMongoRepo, productRedisRepo, s.log, productsProducer)

	im := interceptors.NewInterceptorManager(s.log, s.cfg)
	mw := middlewares.NewMiddlewareManager(s.log, s.cfg)

	l, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	defer l.Close()

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		s.log.Fatalf("failed to load key pair: %s", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: s.cfg.Server.MaxConnectionIdle * time.Minute,
			Timeout:           s.cfg.Server.Timeout * time.Second,
			MaxConnectionAge:  s.cfg.Server.MaxConnectionAge * time.Minute,
			Time:              s.cfg.Server.Timeout * time.Minute,
		}),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
			im.Logger,
		),
	)

	productService := product.NewProductService(s.log, productUC, validate)
	productsService.RegisterProductsServiceServer(grpcServer, productService)
	grpc_prometheus.Register(grpcServer)

	v1 := s.echo.Group("/api/v1")
	v1.Use(mw.Metrics)

	productHandlers := productsHttpV1.NewProductHandlers(s.log, productUC, validate, v1.Group("/products"), mw)
	productHandlers.MapRoutes()

	productsCG := kafka.NewProductsConsumerGroup(s.cfg.Kafka.Brokers, kafkaGroupID, s.log, s.cfg, productUC, validate)
	productsCG.RunConsumers(ctx, cancel)

	go func() {
		s.log.Infof("Server is listening on PORT: %s", s.cfg.Http.Port)
		s.runHttpServer()
	}()

	go func() {
		s.log.Infof("GRPC Server is listening on port: %s", s.cfg.Server.Port)
		s.log.Fatal(grpcServer.Serve(l))
	}()

	if s.cfg.Server.Development {
		reflection.Register(grpcServer)
	}

	metricsServer := echo.New()
	go func() {
		metricsServer.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		s.log.Infof("Metrics server is running on port: %s", s.cfg.Metrics.Port)
		if err := metricsServer.Start(s.cfg.Metrics.Port); err != nil {
			s.log.Error(err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.log.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.log.Errorf("ctx.Done: %v", done)
	}

	if err := s.echo.Server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "echo.Server.Shutdown")
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		s.log.Errorf("metricsServer.Shutdown: %v", err)
	}
	grpcServer.GracefulStop()
	s.log.Info("Server Exited Properly")

	return nil
}
