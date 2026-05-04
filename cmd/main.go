package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/satriaardiperdana-2020/monelog/internal/api"
	"github.com/satriaardiperdana-2020/monelog/internal/config"
	"github.com/satriaardiperdana-2020/monelog/internal/handlers"
	authMiddleware "github.com/satriaardiperdana-2020/monelog/internal/middleware"
	"github.com/satriaardiperdana-2020/monelog/internal/repository/postgresql"
	"log"
)

func main() {
	cfg, err := config.Load("config-development.yml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	pool, err := postgresql.NewConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	queries := postgresql.New(pool)

	authHandler := &handlers.AuthHandler{Queries: queries, JWTSecret: []byte(cfg.JWT.Secret)}
	txHandler := &handlers.TransactionHandler{Queries: queries}
	catHandler := &handlers.CategoryHandler{Queries: queries}
	userHandler := &handlers.UserHandler{Queries: queries}
	server := &handlers.Server{
		Transaction: txHandler,
		Category:    catHandler,
		User:        userHandler,
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public routes (no JWT middleware)
	e.POST("/auth/register", authHandler.Register)
	e.POST("/auth/login", authHandler.Login)

	// Protected routes
	apiGroup := e.Group("/api/v1")
	apiGroup.Use(authMiddleware.JWTAuth([]byte(cfg.JWT.Secret)))
	strictHandler := api.NewStrictHandler(server, nil)
	api.RegisterHandlers(apiGroup, strictHandler)

	log.Fatal(e.Start(":" + cfg.Server.Port))
}
