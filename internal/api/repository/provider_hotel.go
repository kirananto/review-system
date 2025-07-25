package repository

import (
	"github.com/kirananto/review-system/internal/api/dto"
	models "github.com/kirananto/review-system/internal/models"
)

// GetProviderHotelsList retrieves all provider hotels.
func (r *reviewRepository) GetProviderHotelsList(queryParams *dto.ProviderHotelsQueryParams) ([]*models.ProviderHotel, int, error) {
	var providerHotels []*models.ProviderHotel
	var totalCount int64

	// Initialize query
	dbQuery := r.db.Model(&models.ProviderHotel{})

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
		Find(&providerHotels).Error; err != nil {
		return nil, 0, err
	}

	// Get total count using the same conditions
	if err := dbQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	return providerHotels, int(totalCount), nil
}

// GetProviderHotel retrieves a provider-specific hotel mapping.
func (r *reviewRepository) GetProviderHotel(providerID uint, hotelID uint) (*models.ProviderHotel, error) {
	var providerHotel models.ProviderHotel
	if err := r.db.Where("provider_id = ? AND hotel_id = ?", providerID, hotelID).First(&providerHotel).Error; err != nil {
		return nil, err
	}
	return &providerHotel, nil
}

// CreateProviderHotel creates a new provider-specific hotel mapping.
func (r *reviewRepository) CreateProviderHotel(providerHotel *models.ProviderHotel) error {
	return r.db.Create(providerHotel).Error
}

// UpdateProviderHotel updates an existing provider-specific hotel mapping.
func (r *reviewRepository) UpdateProviderHotel(providerHotel *models.ProviderHotel) error {
	return r.db.Save(providerHotel).Error
}
