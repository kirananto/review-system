package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/kirananto/review-system/docs"
	"github.com/kirananto/review-system/internal/config"
	"github.com/kirananto/review-system/internal/logger"
	"github.com/kirananto/review-system/internal/server"
)

// @title Review System API
// @version 1.0
// @description This is a sample server for a review system.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	appCfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create server config
	serverCfg := &server.ServerConfig{
		DatabaseDSN: appCfg.Database.DSN,
		Port:        os.Getenv("PORT"),
		RunMode:     os.Getenv("RUN_MODE"),
		LogConfig: logger.LogConfig{
			LogLevel: os.Getenv("LOG_LEVEL"),
		},
	}

	// Create and start server
	srv, err := server.NewServer(serverCfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
