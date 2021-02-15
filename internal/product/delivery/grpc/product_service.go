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

func (p *productService) Create(ctx context.Context, req *productsService.CreateReq) (*productsService.CreateRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productService.Create")
	defer span.Finish()

	catID, err := primitive.ObjectIDFromHex(req.GetCategoryID())
	if err != nil {
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	prod := &models.Product{
		CategoryID:  catID,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		ImageURL:    req.GetImageURL(),
		Photos:      req.GetPhotos(),
		Quantity:    req.GetQuantity(),
		Rating:      int(req.GetRating()),
	}

	if err := p.validate.StructCtx(ctx, prod); err != nil {
		p.log.Errorf("validate.StructCtx: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	created, err := p.productUC.Create(ctx, prod)
	if err != nil {
		p.log.Errorf("productUC.Create: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	return &productsService.CreateRes{Product: created.ToProto()}, nil
}

func (p *productService) Update(ctx context.Context, req *productsService.UpdateReq) (*productsService.UpdateRes, error) {
	panic("implement me")
}

func (p *productService) GetByID(ctx context.Context, req *productsService.GetByIDReq) (*productsService.GetByIDRes, error) {
	panic("implement me")
}

func (p *productService) Search(ctx context.Context, req *productsService.SearchReq) (*productsService.SearchRes, error) {
	panic("implement me")
}
