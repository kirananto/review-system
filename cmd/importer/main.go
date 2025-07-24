package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kirananto/review-system/internal/api/repository"
	"github.com/kirananto/review-system/internal/api/service"
	"github.com/kirananto/review-system/internal/config"
	"github.com/kirananto/review-system/internal/db"
	"github.com/kirananto/review-system/internal/logger"
	models "github.com/kirananto/review-system/internal/models"
)

func main() {

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dataSource := db.NewDataSource(cfg.Database.DSN)

	dataSource.Db.AutoMigrate(&models.Provider{}, &models.Hotel{}, &models.Review{}, &models.ProviderHotel{})

	log := logger.NewLogger(&logger.LogConfig{LogLevel: "info"})
	repository := repository.NewReviewRepository(dataSource)
	service := service.NewReviewService(repository, log)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file-path>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		log.Error(err, fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()

	if err := service.ProcessReviews(context.Background(), file); err != nil {
		log.Error(err, fmt.Sprintf("Failed to process reviews: %v", err))
	}

	fmt.Println("Successfully processed reviews from", filePath)
}
