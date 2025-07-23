package review

import "time"

// Provider represents a review provider (e.g., Agoda, Booking.com).
type Provider struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"` // e.g., "Agoda"
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
	ID              uint      `json:"id" gorm:"primaryKey"`
	HotelID         uint      `json:"hotel_id" gorm:"not null;index"`
	ProviderID      uint      `json:"provider_id" gorm:"not null;index"`
	ProviderHotelID string    `json:"provider_hotel_id" gorm:"not null"` // Provider's hotelId (e.g., "10984")
	OverallScore    float64   `json:"overall_score" gorm:"default:0"`
	ReviewCount     int       `json:"review_count" gorm:"default:0"`
	Grades          string    `json:"grades" gorm:"type:json"` // JSON-encoded grades
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// Review represents a single review from a provider.
type Review struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	ProviderID       uint      `json:"provider_id" gorm:"not null;index:idx_provider_review,unique"`
	ProviderHotelID  uint      `json:"provider_hotel_id" gorm:"not null;references:ProviderHotel"`
	ProviderReviewID string    `json:"provider_review_id" gorm:"not null;index:idx_provider_review,unique"` // From comment.hotelReviewId
	Rating           float64   `json:"rating" gorm:"not null"`
	Comment          string    `json:"comment"`
	Lang             string    `json:"lang" gorm:"default:'en'"`
	ReviewDate       time.Time `json:"review_date" gorm:"not null"`
	ReviewerInfo     string    `json:"reviewer_info" gorm:"type:json"` // JSON-encoded reviewer info
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
}
