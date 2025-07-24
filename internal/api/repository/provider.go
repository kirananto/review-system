package repository

import (
	"github.com/kirananto/review-system/internal/api/dto"
	models "github.com/kirananto/review-system/internal/models"
)

// GetProvidersList retrieves providers with pagination and filters
func (r *reviewRepository) GetProvidersList(queryParams *dto.ProvidersQueryParams) ([]*models.Provider, int, error) {
	var providers []*models.Provider
	var totalCount int64

	// Initialize query
	dbQuery := r.db.Model(&models.Provider{})

	// Apply filters
	if queryParams.Name != "" {
		dbQuery = dbQuery.Where("name ILIKE ?", "%"+queryParams.Name+"%")
	}

	// Get paginated results
	if err := dbQuery.
		Order("updated_at desc").
		Offset(queryParams.Offset).
		Limit(queryParams.Limit).
		Find(&providers).Error; err != nil {
		return nil, 0, err
	}

	// Get total count using the same conditions
	if err := dbQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	return providers, int(totalCount), nil
}

// GetProviderByID retrieves a provider by its ID.
func (r *reviewRepository) GetProviderByID(id uint) (*models.Provider, error) {
	var provider models.Provider
	if err := r.db.First(&provider, id).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

// GetProviderByName retrieves a provider by its name.
func (r *reviewRepository) GetProviderByName(name string) (*models.Provider, error) {
	var provider models.Provider
	if err := r.db.Where("name = ?", name).First(&provider).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

// CreateProvider creates a new provider.
func (r *reviewRepository) CreateProvider(provider *models.Provider) error {
	return r.db.Create(provider).Error
}
