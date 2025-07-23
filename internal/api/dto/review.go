package dto

import "time"

type ReviewRequest struct {
	ProductID string  `json:"product_id"`
	Rating    float32 `json:"rating"`
	Comment   string  `json:"comment"`
}

type ReviewResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Rating    float32   `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
