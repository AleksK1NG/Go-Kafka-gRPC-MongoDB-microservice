package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/AleksK1NG/products-microservice/config"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
)

var (
	httpTotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_microservice_total_requests",
		Help: "The total number of incoming HTTP requests",
	})
)

// MiddlewareManager http middlewares
type middlewareManager struct {
	log logger.Logger
	cfg *config.Config
}

// MiddlewareManager interface
type MiddlewareManager interface {
	Metrics(next echo.HandlerFunc) echo.HandlerFunc
}

// NewMiddlewareManager constructor
func NewMiddlewareManager(log logger.Logger, cfg *config.Config) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg}
}

// Metrics prometheus metrics
func (m *middlewareManager) Metrics(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		httpTotalRequests.Inc()
		return next(c)
	}
}
