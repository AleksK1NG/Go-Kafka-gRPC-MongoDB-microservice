package grpc

import (
	"github.com/AleksK1NG/products-microservice/internal/product"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
)

// productGRPCService gRPC Service
type productGRPCService struct {
	log       logger.Logger
	productUC product.UseCase
}
