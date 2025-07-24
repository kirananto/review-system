package validator

import (
	"errors"

	"github.com/kirananto/review-system/internal/api/dto"
)

type ReviewValidator struct{}

func NewReviewValidator() *ReviewValidator {
	return &ReviewValidator{}
}

func (v *ReviewValidator) ValidateCreateReview(req *dto.ReviewRequest) error {
	if req.ProductID == "" {
		return errors.New("product_id is required")
	}
	if req.Rating < 0 || req.Rating > 5 {
		return errors.New("rating must be between 0 and 5")
	}

	if req.ExternalID == "" {
		return errors.New("external_id required")
	}
	return nil
}
