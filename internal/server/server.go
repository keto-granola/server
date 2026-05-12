package server

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/keto-granola/server/internal/config"
	"github.com/keto-granola/server/internal/store/db"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

const (
	// How long the server will wait to read the entire request after the connection is accepted
	readTimeout = 10 * time.Second

	// How long the server has to write the response after reading the request
	writeTimeout = 10 * time.Second

	// How long to keep a keep-alive connection open waiting for the next request
	idleTimeout = 120 * time.Second

	shutdownTimeout  = 10 * time.Second
	serverRateLimit  = 60
	serverBurstLimit = 120
)

type Server struct {
	instance *echo.Echo
	port     string
}

type Dependencies struct {
	Db *db.Db
}

type customValidator struct {
	validator *validator.Validate
}

func New(ctx context.Context, deps Dependencies, cfg *config.App) *Server {
	instance := echo.New()
	instance.Validator = &customValidator{validator: validator.New()}
	instance.HideBanner = true // Prevents startup banner from being logged

	instance.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{cfg.ClientURL},
	}))

	// TODO: implement logging middleware
	// instance.Use(middleware.Log)

	// limits each unique IP to 60 requests per minute with a burst of 120.
	instance.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
		echoMiddleware.RateLimiterMemoryStoreConfig{
			Rate:      serverRateLimit,
			Burst:     serverBurstLimit,
			ExpiresIn: time.Minute,
		},
	)))

	handlers := NewHandlers(
		deps.Db,
	)

	if cfg.Environment == config.EnvironmentTest {
		// TODO: run test middleware
	} else {
		// TODO: run auth middleware
	}

	public := instance.Group("")
	private := instance.Group("")

	registerRoutes(public, private, handlers)

	instance.Server.ReadTimeout = readTimeout
	instance.Server.WriteTimeout = writeTimeout
	instance.Server.IdleTimeout = idleTimeout

	return &Server{
		instance: instance,
		port:     cfg.Port,
	}
}

func (cv *customValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func (s *Server) Start() {

}

func (s *Server) Stop() {

}
