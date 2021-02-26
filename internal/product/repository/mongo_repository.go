package repository

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/AleksK1NG/products-microservice/internal/models"
	productErrors "github.com/AleksK1NG/products-microservice/pkg/product_errors"
	"github.com/AleksK1NG/products-microservice/pkg/utils"
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

	return product, nil
}

// Update Single product
func (p *productMongoRepo) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productMongoRepo.Update")
	defer span.Finish()

	collection := p.mongoDB.Database(productsDB).Collection(productsCollection)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var prod models.Product
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": product.ProductID}, bson.M{"$set": product}, ops).Decode(&prod); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

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

// Search Search product
func (p *productMongoRepo) Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.ProductsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productMongoRepo.Search")
	defer span.Finish()

	collection := p.mongoDB.Database(productsDB).Collection(productsCollection)

	f := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: primitive.Regex{
				Pattern: search,
				Options: "gi",
			}}},
			bson.D{{Key: "description", Value: primitive.Regex{
				Pattern: search,
				Options: "gi",
			}}},
		}},
	}

	count, err := collection.CountDocuments(ctx, f)
	if err != nil {
		return nil, errors.Wrap(err, "CountDocuments")
	}
	if count == 0 {
		return &models.ProductsList{
			TotalCount: 0,
			TotalPages: 0,
			Page:       0,
			Size:       0,
			HasMore:    false,
			Products:   make([]*models.Product, 0),
		}, nil
	}

	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	cursor, err := collection.Find(ctx, f, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx)

	products := make([]*models.Product, 0, pagination.GetSize())
	for cursor.Next(ctx) {
		var prod models.Product
		if err := cursor.Decode(&prod); err != nil {
			return nil, errors.Wrap(err, "Find")
		}
		products = append(products, &prod)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor.Err")
	}

	return &models.ProductsList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Products:   products,
	}, nil
}
