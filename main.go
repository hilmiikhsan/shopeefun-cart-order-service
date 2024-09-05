package main

import (
	"database/sql"
	"log"

	"github.com/hilmiikhsan/shopeefun-cart-order-service/config"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/routes"
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
	return &routes.Routes{}
}
