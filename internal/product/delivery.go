package product

import "github.com/labstack/echo/v4"

// HttpDelivery http delivery
type HttpDelivery interface {
	CreateProduct() echo.HandlerFunc
	UpdateProduct() echo.HandlerFunc
	GetByIDProduct() echo.HandlerFunc
	SearchProduct() echo.HandlerFunc
}
