package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/AleksK1NG/products-microservice/internal/models"
	"github.com/AleksK1NG/products-microservice/internal/product"
	grpcErrors "github.com/AleksK1NG/products-microservice/pkg/grpc_errors"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
	"github.com/AleksK1NG/products-microservice/pkg/utils"
	"github.com/AleksK1NG/products-microservice/proto/product"
)

// productService gRPC Service
type productService struct {
	log       logger.Logger
	productUC product.UseCase
	validate  *validator.Validate
}

// NewProductService productService constructor
func NewProductService(log logger.Logger, productUC product.UseCase, validate *validator.Validate) *productService {
	return &productService{log: log, productUC: productUC, validate: validate}
}

// Create create new product
func (p *productService) Create(ctx context.Context, req *productsService.CreateReq) (*productsService.CreateRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productService.Create")
	defer span.Finish()
	createMessages.Inc()

	catID, err := primitive.ObjectIDFromHex(req.GetCategoryID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	prod := &models.Product{
		CategoryID:  catID,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		ImageURL:    &req.ImageURL,
		Photos:      req.GetPhotos(),
		Quantity:    req.GetQuantity(),
		Rating:      int(req.GetRating()),
	}

	created, err := p.productUC.Create(ctx, prod)
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("productUC.Create: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	successMessages.Inc()
	return &productsService.CreateRes{Product: created.ToProto()}, nil
}

// Update Update existing product
func (p *productService) Update(ctx context.Context, req *productsService.UpdateReq) (*productsService.UpdateRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productService.Update")
	defer span.Finish()
	updateMessages.Inc()

	prodID, err := primitive.ObjectIDFromHex(req.GetProductID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}
	catID, err := primitive.ObjectIDFromHex(req.GetCategoryID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	prod := &models.Product{
		ProductID:   prodID,
		CategoryID:  catID,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		ImageURL:    &req.ImageURL,
		Photos:      req.GetPhotos(),
		Quantity:    req.GetQuantity(),
		Rating:      int(req.GetRating()),
	}

	update, err := p.productUC.Update(ctx, prod)
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("productUC.Update: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	successMessages.Inc()
	return &productsService.UpdateRes{Product: update.ToProto()}, nil
}

// GetByID Get single product by id
func (p *productService) GetByID(ctx context.Context, req *productsService.GetByIDReq) (*productsService.GetByIDRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productService.GetByID")
	defer span.Finish()
	getByIdMessages.Inc()

	prodID, err := primitive.ObjectIDFromHex(req.GetProductID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	prod, err := p.productUC.GetByID(ctx, prodID)
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("productUC.GetByID: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	successMessages.Inc()
	return &productsService.GetByIDRes{Product: prod.ToProto()}, nil
}

// Search Search products
func (p *productService) Search(ctx context.Context, req *productsService.SearchReq) (*productsService.SearchRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productService.Search")
	defer span.Finish()
	searchMessages.Inc()

	products, err := p.productUC.Search(ctx, req.GetSearch(), utils.NewPaginationQuery(int(req.GetSize()), int(req.GetPage())))
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("productUC.Search: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	successMessages.Inc()
	return &productsService.SearchRes{
		TotalCount: products.TotalCount,
		TotalPages: products.TotalPages,
		Page:       products.Page,
		Size:       products.Size,
		HasMore:    products.HasMore,
		Products:   products.ToProtoList(),
	}, nil
}
