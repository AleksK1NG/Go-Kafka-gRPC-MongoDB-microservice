package repository

import (
	"context"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/AleksK1NG/products-microservice/internal/models"
	productErrors "github.com/AleksK1NG/products-microservice/pkg/product_errors"
)

const (
	productsDB         = "products"
	productsCollection = "products"
)

// productMongoRepo
type productMongoRepo struct {
	mongoDB *mongo.Client
}

// NewProductMongoRepo productMongoRepo constructor
func NewProductMongoRepo(mongoDB *mongo.Client) *productMongoRepo {
	return &productMongoRepo{mongoDB: mongoDB}
}

// Create Create new product
func (p *productMongoRepo) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productMongoRepo.Create")
	defer span.Finish()

	collection := p.mongoDB.Database(productsDB).Collection(productsCollection)

	product.CreatedAt = time.Now().UTC()
	product.UpdatedAt = time.Now().UTC()

	result, err := collection.InsertOne(ctx, product, &options.InsertOneOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "InsertOne")
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.Wrap(productErrors.ErrObjectIDTypeConversion, "result.InsertedID")
	}

	product.ProductID = objectID

	log.Printf("CREATED PRODUCT: %-v", product)

	return product, nil
}

// Update Single product
func (p *productMongoRepo) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productMongoRepo.Update")
	defer span.Finish()

	collection := p.mongoDB.Database(productsDB).Collection(productsCollection)

	upsert := true
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	product.UpdatedAt = time.Now().UTC()
	var prod models.Product
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": product.ProductID}, bson.M{"$set": product}, opts).Decode(&prod); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	log.Printf("UPDATED PRODUCT: %-v", product)

	return &prod, nil
}

// GetByID Get single product by id
func (p *productMongoRepo) GetByID(ctx context.Context, productID primitive.ObjectID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productMongoRepo.GetByID")
	defer span.Finish()

	collection := p.mongoDB.Database(productsDB).Collection(productsCollection)

	var prod models.Product
	if err := collection.FindOne(ctx, bson.M{"_id": productID}).Decode(&prod); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &prod, nil
}

func (p *productMongoRepo) Search(ctx context.Context, search string, page, size int64) ([]*models.Product, error) {
	panic("implement me")
}
