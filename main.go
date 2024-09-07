package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/hilmiikhsan/shopeefun-cart-order-service/config"
	cartHandler "github.com/hilmiikhsan/shopeefun-cart-order-service/handlers/cart"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/repository/cart"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/routes"
	cartUsecase "github.com/hilmiikhsan/shopeefun-cart-order-service/usecase/cart"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/validators"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
		return
	}

	sqlDb, err := config.ConnectToDatabase(config.Connection{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
		return
	}

	routes := setupRoutes(sqlDb)
	routes.Run(cfg.AppPort)
}

func setupRoutes(db *sql.DB) *routes.Routes {
	ctx := context.Background()
	validatorInstance := validators.NewValidator()

	cartRepository := cart.NewStore(db)
	cartUseCase := cartUsecase.NewCart(ctx, cartRepository)
	cartHandler := cartHandler.NewHandler(cartUseCase, validatorInstance)

	return &routes.Routes{
		Cart: cartHandler,
	}
}
