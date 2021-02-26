package product

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/AleksK1NG/products-microservice/internal/models"
	"github.com/AleksK1NG/products-microservice/pkg/utils"
)

// UseCase Product
type UseCase interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	GetByID(ctx context.Context, productID primitive.ObjectID) (*models.Product, error)
	Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.ProductsList, error)
	PublishCreate(ctx context.Context, product *models.Product) error
	PublishUpdate(ctx context.Context, product *models.Product) error
}
