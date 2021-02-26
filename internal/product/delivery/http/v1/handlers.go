package v1

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/AleksK1NG/products-microservice/internal/middlewares"
	"github.com/AleksK1NG/products-microservice/internal/models"
	"github.com/AleksK1NG/products-microservice/internal/product"
	httpErrors "github.com/AleksK1NG/products-microservice/pkg/http_errors"

	"github.com/AleksK1NG/products-microservice/pkg/logger"
	"github.com/AleksK1NG/products-microservice/pkg/utils"
)

type productHandlers struct {
	log       logger.Logger
	productUC product.UseCase
	validate  *validator.Validate
	group     *echo.Group
	mw        middlewares.MiddlewareManager
}

// NewProductHandlers constructor
func NewProductHandlers(
	log logger.Logger,
	productUC product.UseCase,
	validate *validator.Validate,
	group *echo.Group,
	mw middlewares.MiddlewareManager,
) *productHandlers {
	return &productHandlers{log: log, productUC: productUC, validate: validate, group: group, mw: mw}
}

// CreateProduct Create product
// @Tags Products
// @Summary Create new product
// @Description Create new single product
// @Accept json
// @Produce json
// @Success 201 {object} models.Product
// @Router /products [post]
func (p *productHandlers) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "productHandlers.Create")
		defer span.Finish()
		createRequests.Inc()

		var prod models.Product
		if err := c.Bind(&prod); err != nil {
			p.log.Errorf("c.Bind: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := p.validate.StructCtx(ctx, &prod); err != nil {
			p.log.Errorf("validate.StructCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := p.productUC.PublishCreate(ctx, &prod); err != nil {
			p.log.Errorf("productUC.PublishCreate: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.NoContent(http.StatusCreated)
	}
}

// UpdateProduct Update product
// @Tags Products
// @Summary Update single product
// @Description Update single product by id
// @Accept json
// @Produce json
// @Param product_id path string true "product id"
// @Success 200 {object} models.Product
// @Router /products/{product_id} [put]
func (p *productHandlers) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "productHandlers.Update")
		defer span.Finish()
		updateRequests.Inc()

		var prod models.Product
		if err := c.Bind(&prod); err != nil {
			p.log.Errorf("c.Bind: %v", err)
			errorRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err)
		}

		prodID, err := primitive.ObjectIDFromHex(c.Param("product_id"))
		if err != nil {
			p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
			errorRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err)
		}
		prod.ProductID = prodID

		if err := p.validate.StructCtx(ctx, &prod); err != nil {
			p.log.Errorf("validate.StructCtx: %v", err)
			errorRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := p.productUC.PublishUpdate(ctx, &prod); err != nil {
			p.log.Errorf("productUC.PublishUpdate: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.NoContent(http.StatusOK)
	}
}

// GetByIDProduct Get product by id
// @Tags Products
// @Summary Get product by id
// @Description Get single product by id
// @Accept json
// @Produce json
// @Param product_id path string true "product id"
// @Success 200 {object} models.Product
// @Router /products/{product_id} [get]
func (p *productHandlers) GetByIDProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "productHandlers.GetByID")
		defer span.Finish()
		getByIdRequests.Inc()

		objectID, err := primitive.ObjectIDFromHex(c.Param("product_id"))
		if err != nil {
			p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
			errorRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err)
		}

		prod, err := p.productUC.GetByID(ctx, objectID)
		if err != nil {
			p.log.Errorf("productUC.GetByID: %v", err)
			errorRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.JSON(http.StatusOK, prod)
	}
}

// SearchProduct Search product
// @Tags Products
// @Summary Search product
// @Description Search product by name or description
// @Accept json
// @Produce json
// @Param search query string false "search text"
// @Param page query string false "page number"
// @Param size query string false "number of elements"
// @Success 200 {object} models.ProductsList
// @Router /products/search [get]
func (p *productHandlers) SearchProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "productHandlers.Search")
		defer span.Finish()
		searchRequests.Inc()

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			p.log.Errorf("strconv.Atoi: %v", err)
			errorRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest)
		}
		size, err := strconv.Atoi(c.QueryParam("size"))
		if err != nil {
			p.log.Errorf("strconv.Atoi: %v", err)
			errorRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest)
		}

		pq := utils.NewPaginationQuery(size, page)
		result, err := p.productUC.Search(ctx, c.QueryParam("search"), pq)
		if err != nil {
			p.log.Errorf("productUC.Search: %v", err)
			errorRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.JSON(http.StatusOK, result)
	}
}
