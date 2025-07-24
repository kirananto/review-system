package repository

import (
	"github.com/kirananto/review-system/internal/api/dto"
	reviewmodel "github.com/kirananto/review-system/pkg/review"
)

func (r *reviewRepository) GetReviewsList(queryParams *dto.ReviewQueryParams) ([]*reviewmodel.Review, int, error) {
	var reviews []*reviewmodel.Review
	var totalCount int64

	// Initialize query
	dbQuery := r.db.Model(&reviewmodel.Review{})

	// Build conditions map with only non-zero values
	conditions := make(map[string]interface{})
	if queryParams.HotelID != 0 {
		conditions["hotel_id"] = queryParams.HotelID
	}
	if queryParams.ProviderID != 0 {
		conditions["provider_id"] = queryParams.ProviderID
	}

	// Apply non-zero conditions (GORM will AND them together)
	if len(conditions) > 0 {
		dbQuery = dbQuery.Where(conditions)
	}

	// Get paginated results
	if err := dbQuery.
		Order("updated_at desc").
		Offset(queryParams.Offset).
		Limit(queryParams.Limit).
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	// Get total count using the same conditions
	if err := dbQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	return reviews, int(totalCount), nil
}

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

// GetReviews retrieves all reviews.
// TODO: Use GetReviewsList with pagination and filters instead of this method.
func (r *reviewRepository) GetReviews() ([]*reviewmodel.Review, error) {
	var reviews []*reviewmodel.Review
	if err := r.db.Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}
