package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	authMiddleware "github.com/satriaardiperdana-2020/monelog-be/internal/middleware"
	"log"
	"net/http"

	"github.com/satriaardiperdana-2020/monelog-be/internal/api"
	"github.com/satriaardiperdana-2020/monelog-be/internal/config"
	"github.com/satriaardiperdana-2020/monelog-be/internal/handlers"
	"github.com/satriaardiperdana-2020/monelog-be/internal/repository/postgresql"
)

func main() {
	// Load configuration from YAML file
	cfg, err := config.Load("config-development.yml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to PostgreSQL database
	pool, err := postgresql.NewConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	queries := postgresql.New(pool)

	// Initialize all handlers
	authHandler := &handlers.AuthHandler{
		Queries:   queries,
		JWTSecret: []byte(cfg.JWT.Secret),
	}
	txHandler := &handlers.TransactionHandler{Queries: queries}
	catHandler := &handlers.CategoryHandler{Queries: queries}
	userHandler := &handlers.UserHandler{Queries: queries}

	// Combine handlers into a single server that implements the strict interface
	server := &handlers.Server{
		Queries: queries,
		//JWTSecret:   []byte(cfg.JWT.Secret),
		Auth:        authHandler,
		Category:    catHandler,
		Transaction: txHandler,
		User:        userHandler,
	}

	// Create Echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// CORS configuration – allows frontend origin (adjust as needed)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost", "http://10.222.136.79:8080"}, // Vue dev server
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	apiGroup := e.Group("/api/v1")
	apiGroup.Use(authMiddleware.JWTAuth([]byte(cfg.JWT.Secret), queries))
	apiGroup.POST("/auth/logout", authHandler.Logout)
	//apiGroup.Use(authMiddleware.JWTAuth([]byte(cfg.JWT.Secret)))
	strictHandler := api.NewStrictHandler(server, nil) // second argument = error handler (optional)
	api.RegisterHandlers(apiGroup, strictHandler)
	// ----- Conditional JWT middleware on the same group -----
	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			// Allow registration and login without token
			if path == "/api/v1/auth/register" || path == "/api/v1/auth/login" {
				return next(c)
			}
			// All other endpoints under /api/v1 need a valid JWT
			return authMiddleware.JWTAuth([]byte(cfg.JWT.Secret), queries)(next)(c)
		}
	})

	// Start server
	log.Fatal(e.Start(":" + cfg.Server.Port))
}
