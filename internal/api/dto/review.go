package dto

import "time"

type ReviewRequest struct {
	ProductID  string  `json:"product_id"`
	Rating     float32 `json:"rating"`
	ExternalID string  `json:"external_id"`
	Comment    string  `json:"comment"`
}

type ReviewResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Rating    float32   `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewQueryParams struct {
	Limit      int  `schema:"limit"`
	Offset     int  `schema:"offset"`
	HotelID    uint `schema:"hotel_id"`
	ProviderID uint `schema:"provider_id"`
}
