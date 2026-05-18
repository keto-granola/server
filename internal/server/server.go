package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/keto-granola/server/internal/config"
	"github.com/keto-granola/server/internal/middleware"
	productadmin "github.com/keto-granola/server/internal/product/admin"
	"github.com/keto-granola/server/internal/store"
)

const (
	// how long the server will wait to read the entire request after the connection is accepted
	readTimeout = 10 * time.Second
	// how long the server has to write the response after reading the request
	writeTimeout = 10 * time.Second
	// how long to keep a keep-alive connection open waiting for the next request
	idleTimeout     = 120 * time.Second
	shutdownTimeout = 10 * time.Second

	serverRateLimit  = 60
	serverBurstLimit = 120
)

type Server struct {
	instance *echo.Echo
}

type Handlers struct {
	ProductAdmin *productadmin.Handler
}

type customValidator struct {
	validator *validator.Validate
}

func New(ctx context.Context, environment config.Environment, clientURL string, handlers *Handlers, dataStore *store.Store) *Server {
	instance := echo.New()
	instance.Validator = &customValidator{validator: validator.New()}
	instance.HideBanner = true // Prevents startup banner from being logged

	instance.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{clientURL},
	}))

	instance.Use(middleware.Log)

	// limits each unique IP to 60 requests per minute with a burst of 120.
	instance.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
		echoMiddleware.RateLimiterMemoryStoreConfig{
			Rate:      serverRateLimit,
			Burst:     serverBurstLimit,
			ExpiresIn: time.Minute,
		},
	)))

	if environment == config.EnvironmentTest {
		// TODO: run test middleware
		slog.Info("run test middleware")
	} else {
		// TODO: run auth middleware
		slog.Info("run auth middleware")
	}

	public := instance.Group("/" + config.APIVersion)
	private := instance.Group("/" + config.APIVersion)

	registerRoutes(public, private, handlers, dataStore)

	instance.Server.ReadTimeout = readTimeout
	instance.Server.WriteTimeout = writeTimeout
	instance.Server.IdleTimeout = idleTimeout

	return &Server{
		instance: instance,
	}
}

func (cv *customValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func (s *Server) Start(port string) error {
	if err := s.instance.Start(":" + port); err != nil && err != http.ErrServerClosed {
		slog.Error("start server", slog.Any("error", err))
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return s.instance.Shutdown(ctx)
}
