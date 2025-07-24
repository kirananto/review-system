package repository

import (
	"github.com/kirananto/review-system/internal/api/dto"
	models "github.com/kirananto/review-system/internal/models"
)

func (r *reviewRepository) GetReviewsList(queryParams *dto.ReviewQueryParams) ([]*models.Review, int, error) {
	var reviews []*models.Review
	var totalCount int64

	// Initialize query
	dbQuery := r.db.Model(&models.Review{})

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
func (r *reviewRepository) GetReviewByID(id uint) (*models.Review, error) {
	var review models.Review
	if err := r.db.First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// CreateReview creates a new review.
func (r *reviewRepository) CreateReview(review *models.Review) error {
	return r.db.Create(review).Error
}

// GetReviews retrieves all reviews.
// TODO: Use GetReviewsList with pagination and filters instead of this method.
func (r *reviewRepository) GetReviews() ([]*models.Review, error) {
	var reviews []*models.Review
	if err := r.db.Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}
