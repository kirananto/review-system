package repository

import (
	"github.com/kirananto/review-system/internal/api/dto"
	reviewmodel "github.com/kirananto/review-system/pkg/review"
)

// GetHotelsList retrieves hotels with pagination and filters
func (r *reviewRepository) GetHotelsList(queryParams *dto.HotelsQueryParams) ([]*reviewmodel.Hotel, int, error) {
	var hotels []*reviewmodel.Hotel
	var totalCount int64

	// Initialize query
	dbQuery := r.db.Model(&reviewmodel.Hotel{})

	// Apply filters
	if queryParams.Name != "" {
		dbQuery = dbQuery.Where("hotel_name ILIKE ?", "%"+queryParams.Name+"%")
	}

	// Get paginated results
	if err := dbQuery.
		Order("updated_at desc").
		Offset(queryParams.Offset).
		Limit(queryParams.Limit).
		Find(&hotels).Error; err != nil {
		return nil, 0, err
	}

	// Get total count using the same conditions
	if err := dbQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	return hotels, int(totalCount), nil
}

// GetHotelByID retrieves a hotel by its ID.
func (r *reviewRepository) GetHotelByID(id uint) (*reviewmodel.Hotel, error) {
	var hotel reviewmodel.Hotel
	if err := r.db.First(&hotel, id).Error; err != nil {
		return nil, err
	}
	return &hotel, nil
}

// GetHotelByName retrieves a hotel by its name.
func (r *reviewRepository) GetHotelByName(name string) (*reviewmodel.Hotel, error) {
	var hotel reviewmodel.Hotel
	if err := r.db.Where("hotel_name = ?", name).First(&hotel).Error; err != nil {
		return nil, err
	}
	return &hotel, nil
}

// CreateHotel creates a new hotel.
func (r *reviewRepository) CreateHotel(hotel *reviewmodel.Hotel) error {
	return r.db.Create(hotel).Error
}

// UpdateHotel updates an existing hotel.
func (r *reviewRepository) UpdateHotel(hotel *reviewmodel.Hotel) error {
	return r.db.Save(hotel).Error
}

// DeleteHotel deletes a hotel by its ID.
func (r *reviewRepository) DeleteHotel(id uint) error {
	return r.db.Delete(&reviewmodel.Hotel{}, id).Error
}
