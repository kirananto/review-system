package repository

import (
	"github.com/kirananto/review-system/internal/api/dto"
	"github.com/kirananto/review-system/internal/db"
	models "github.com/kirananto/review-system/internal/models"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	// Provider methods
	GetProvidersList(queryParams *dto.ProvidersQueryParams) ([]*models.Provider, int, error)
	GetProviderByID(id uint) (*models.Provider, error)
	GetProviderByName(name string) (*models.Provider, error)
	CreateProvider(provider *models.Provider) error

	// Hotel methods
	GetHotelsList(queryParams *dto.HotelsQueryParams) ([]*models.Hotel, int, error)
	GetHotelByID(id uint) (*models.Hotel, error)
	GetHotelByName(name string) (*models.Hotel, error)
	CreateHotel(hotel *models.Hotel) error

	// ProviderHotel methods
	GetProviderHotelsList(queryParams *dto.ProviderHotelsQueryParams) ([]*models.ProviderHotel, int, error)
	GetProviderHotel(providerID uint, hotelID uint) (*models.ProviderHotel, error)
	CreateProviderHotel(providerHotel *models.ProviderHotel) error
	UpdateProviderHotel(providerHotel *models.ProviderHotel) error

	// Review methods
	GetReviewsList(queryParams *dto.ReviewQueryParams) ([]*models.Review, int, error)
	GetReviewByID(id uint) (*models.Review, error)
	CreateReview(review *models.Review) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(dataSource *db.DataSource) ReviewRepository {
	return &reviewRepository{
		db: dataSource.Db,
	}
}
