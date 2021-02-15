package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/products-microservice/internal/models"
	"github.com/AleksK1NG/products-microservice/internal/product"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
)

// productUC
type productUC struct {
	productRepo product.MongoRepository
	log         logger.Logger
}

// NewProductUC productUC constructor
func NewProductUC(productRepo product.MongoRepository, log logger.Logger) *productUC {
	return &productUC{productRepo: productRepo, log: log}
}

// Create Create new product
func (p *productUC) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.Create")
	defer span.Finish()
	return p.productRepo.Create(ctx, product)
}

func (p *productUC) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	panic("implement me")
}

func (p *productUC) GetByID(ctx context.Context, productID string) (*models.Product, error) {
	panic("implement me")
}

func (p *productUC) Search(ctx context.Context, search string, page, size int64) ([]*models.Product, error) {
	panic("implement me")
}
