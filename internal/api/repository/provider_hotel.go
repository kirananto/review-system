package repository

import reviewmodel "github.com/kirananto/review-system/pkg/review"

// GetProviderHotels retrieves all provider hotels.
func (r *reviewRepository) GetProviderHotels() ([]*reviewmodel.ProviderHotel, error) {
	var providerHotels []*reviewmodel.ProviderHotel
	if err := r.db.Find(&providerHotels).Error; err != nil {
		return nil, err
	}
	return providerHotels, nil
}

// GetProviderHotelByID retrieves a provider hotel by its ID.
func (r *reviewRepository) GetProviderHotelByID(id uint) (*reviewmodel.ProviderHotel, error) {
	var providerHotel reviewmodel.ProviderHotel
	if err := r.db.First(&providerHotel, id).Error; err != nil {
		return nil, err
	}
	return &providerHotel, nil
}

// GetProviderHotel retrieves a provider-specific hotel mapping.
func (r *reviewRepository) GetProviderHotel(providerID uint, providerHotelID string) (*reviewmodel.ProviderHotel, error) {
	var providerHotel reviewmodel.ProviderHotel
	if err := r.db.Where("provider_id = ? AND provider_hotel_id = ?", providerID, providerHotelID).First(&providerHotel).Error; err != nil {
		return nil, err
	}
	return &providerHotel, nil
}

// CreateProviderHotel creates a new provider-specific hotel mapping.
func (r *reviewRepository) CreateProviderHotel(providerHotel *reviewmodel.ProviderHotel) error {
	return r.db.Create(providerHotel).Error
}

// UpdateProviderHotel updates an existing provider-specific hotel mapping.
func (r *reviewRepository) UpdateProviderHotel(providerHotel *reviewmodel.ProviderHotel) error {
	return r.db.Save(providerHotel).Error
}

// DeleteProviderHotel deletes a provider hotel by its ID.
func (r *reviewRepository) DeleteProviderHotel(id uint) error {
	return r.db.Delete(&reviewmodel.ProviderHotel{}, id).Error
}
