package repository

import (
	reviewmodel "github.com/kirananto/review-system/pkg/review"
)

// GetReviewByID retrieves a review by its ID.
func (r *reviewRepository) GetReviewByID(id uint) (*reviewmodel.Review, error) {
	var review reviewmodel.Review
	if err := r.db.First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// GetReview retrieves a review by its provider and provider-specific review ID.
func (r *reviewRepository) GetReview(providerID uint, providerReviewID string) (*reviewmodel.Review, error) {
	var review reviewmodel.Review
	if err := r.db.Where("provider_id = ? AND provider_review_id = ?", providerID, providerReviewID).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// CreateReview creates a new review.
func (r *reviewRepository) CreateReview(review *reviewmodel.Review) error {
	return r.db.Create(review).Error
}

// UpdateReview updates an existing review.
func (r *reviewRepository) UpdateReview(review *reviewmodel.Review) error {
	return r.db.Save(review).Error
}

// DeleteReview deletes a review by its ID.
func (r *reviewRepository) DeleteReview(id uint) error {
	return r.db.Delete(&reviewmodel.Review{}, id).Error
}

// GetReviews retrieves all reviews.
func (r *reviewRepository) GetReviews() ([]*reviewmodel.Review, error) {
	var reviews []*reviewmodel.Review
	if err := r.db.Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}
