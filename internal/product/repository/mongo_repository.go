package repository

import (
	"context"

	"github.com/AleksK1NG/products-microservice/internal/models"
)

// productMongoRepo
type productMongoRepo struct {
}

func (p *productMongoRepo) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	panic("implement me")
}

func (p *productMongoRepo) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	panic("implement me")
}

func (p *productMongoRepo) GetByID(ctx context.Context, productID string) (*models.Product, error) {
	panic("implement me")
}

func (p *productMongoRepo) Search(ctx context.Context, search string, page, size int64) ([]*models.Product, error) {
	panic("implement me")
}
