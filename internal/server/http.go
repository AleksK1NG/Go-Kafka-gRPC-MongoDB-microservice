package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *server) runHttpServer() {
	s.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})
	s.echo.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	s.mapRoutes()

	go func() {
		s.log.Infof("Server is listening on PORT: %s", s.cfg.Http.Port)
		s.echo.Server.ReadTimeout = time.Second * s.cfg.Http.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.Http.WriteTimeout
		s.echo.Server.MaxHeaderBytes = maxHeaderBytes
		if err := s.echo.StartTLS(s.cfg.Http.Port, certFile, keyFile); err != nil {
			s.log.Fatalf("Error starting TLS Server: ", err)
		}
	}()
}

func (s *server) mapRoutes() {
	// s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	s.echo.Use(middleware.Logger())
	s.echo.Pre(middleware.HTTPSRedirect())
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrfTokenHeader},
	}))
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	s.echo.Use(middleware.Secure())
	s.echo.Use(middleware.BodyLimit(bodyLimit))
}
