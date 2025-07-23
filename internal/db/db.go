package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataSource struct {
	Db *gorm.DB
}

// NewDataSource creates a new database connection.
func NewDataSource(dsn string) *DataSource {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Get the underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get underlying *sql.DB: %v", err))
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Migrate the database schema

	dataSource := &DataSource{Db: db}
	return dataSource
}
