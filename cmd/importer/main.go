package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kirananto/review-system/internal/api/repository"
	"github.com/kirananto/review-system/internal/api/service"
	"github.com/kirananto/review-system/internal/config"
	"github.com/kirananto/review-system/internal/db"
	"github.com/kirananto/review-system/internal/logger"
	reviewmodel "github.com/kirananto/review-system/pkg/review"
)

func main() {
	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dataSource := db.NewDataSource(cfg.Database.DSN)

	// Auto-migrate the schema
	dataSource.Db.AutoMigrate(&reviewmodel.Provider{}, &reviewmodel.Hotel{}, &reviewmodel.ProviderHotel{}, &reviewmodel.Review{})

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
