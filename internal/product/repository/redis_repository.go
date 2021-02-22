package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/AleksK1NG/products-microservice/internal/models"
)

const (
	prefix     = "products"
	expiration = time.Second * 3600
)

type productRedisRepository struct {
	prefix string
	redis  *redis.Client
}

// NewProductRedisRepository constructor
func NewProductRedisRepository(redis *redis.Client) *productRedisRepository {
	return &productRedisRepository{redis: redis, prefix: prefix}
}

func (p *productRedisRepository) SetProduct(ctx context.Context, product *models.Product) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRedisRepository.SetProduct")
	defer span.Finish()

	prodBytes, err := json.Marshal(product)
	if err != nil {
		return errors.Wrap(err, "productRedisRepository.Marshal")
	}

	return p.redis.SetEX(ctx, p.createKey(product.ProductID), string(prodBytes), expiration).Err()
}

func (p *productRedisRepository) GetProductByID(ctx context.Context, productID primitive.ObjectID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRedisRepository.GetProductByID")
	defer span.Finish()

	result, err := p.redis.Get(ctx, p.createKey(productID)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "productRedisRepository.redis.Get")
	}

	var res models.Product
	if err := json.Unmarshal(result, &res); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return &res, nil
}

func (p *productRedisRepository) DeleteProduct(ctx context.Context, productID primitive.ObjectID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRedisRepository.DeleteProduct")
	defer span.Finish()

	return p.redis.Del(ctx, p.createKey(productID)).Err()
}

func (p *productRedisRepository) createKey(id primitive.ObjectID) string {
	return fmt.Sprintf("%s: %s", p.prefix, id.String())
}
