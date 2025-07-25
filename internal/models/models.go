package models

import (
	"encoding/json"
	"time"
)

// Provider represents a review provider (e.g., Agoda, Booking.com).
type Provider struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// Hotel represents a hotel entity.
type Hotel struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	HotelName string    `json:"name" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// ProviderHotel maps a provider's hotel ID to our internal hotel ID.
// It also stores provider-specific overall stats for the hotel.
type ProviderHotel struct {
	HotelID      uint            `json:"hotel_id" gorm:"primaryKey;autoIncrement:false;index:idx_provider_hotel,unique"`
	ProviderID   uint            `json:"provider_id" gorm:"primaryKey;autoIncrement:false;index:idx_provider_hotel,unique"`
	OverallScore float64         `json:"overall_score" gorm:"default:0"`
	ReviewCount  int             `json:"review_count" gorm:"default:0"`
	Grades       json.RawMessage `json:"grades" gorm:"type:jsonb" swaggertype:"string"` // jsonb for Postgres
	UpdatedAt    time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt    time.Time       `json:"created_at" gorm:"autoCreateTime"`

	// enforce FK + cascade to avoid orphans
	Hotel    Hotel    `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignKey:HotelID;references:ID"`
	Provider Provider `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignKey:ProviderID;references:ID"`
}

// Review represents a single review from a provider.
type Review struct {
	ID           uint            `json:"id" gorm:"primaryKey;autoIncrement:false"`
	ProviderID   uint            `json:"provider_id" gorm:"not null"`
	HotelID      uint            `json:"hotel_id" gorm:"not null"`
	Rating       float64         `json:"rating" gorm:"not null"`
	Comment      string          `json:"comment"`
	Lang         string          `json:"lang" gorm:"default:'en'"`
	ReviewDate   time.Time       `json:"review_date" gorm:"not null;index"`
	ReviewerInfo json.RawMessage `json:"reviewer_info" gorm:"type:jsonb" swaggertype:"string"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt    time.Time       `json:"created_at" gorm:"autoCreateTime"`

	Provider Provider `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignKey:ProviderID;references:ID"`
	Hotel    Hotel    `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignKey:HotelID;references:ID"`
}

// AuditLog represents the audit log for a processed file.
type AuditLog struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FileName     string    `json:"file_name" gorm:"not null"`
	SuccessCount int       `json:"success_count" gorm:"default:0"`
	FailureCount int       `json:"failure_count" gorm:"default:0"`
	TotalCount   int       `json:"total_count" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}
