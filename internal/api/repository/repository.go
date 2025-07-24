package repository

import (
	"github.com/kirananto/review-system/internal/api/dto"
	"github.com/kirananto/review-system/internal/db"
	reviewmodel "github.com/kirananto/review-system/pkg/review"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	// Provider methods
	GetProvidersList(queryParams *dto.ProvidersQueryParams) ([]*reviewmodel.Provider, int, error)
	GetProviderByID(id uint) (*reviewmodel.Provider, error)
	GetProviderByName(name string) (*reviewmodel.Provider, error)
	CreateProvider(provider *reviewmodel.Provider) error

	// Hotel methods
	GetHotelsList(queryParams *dto.HotelsQueryParams) ([]*reviewmodel.Hotel, int, error)
	GetHotelByID(id uint) (*reviewmodel.Hotel, error)
	GetHotelByName(name string) (*reviewmodel.Hotel, error)
	CreateHotel(hotel *reviewmodel.Hotel) error

	// ProviderHotel methods
	GetProviderHotelsList(queryParams *dto.ProviderHotelsQueryParams) ([]*reviewmodel.ProviderHotel, int, error)
	GetProviderHotel(providerID uint, hotelID uint) (*reviewmodel.ProviderHotel, error)
	CreateProviderHotel(providerHotel *reviewmodel.ProviderHotel) error
	UpdateProviderHotel(providerHotel *reviewmodel.ProviderHotel) error

	// Review methods
	GetReviewsList(queryParams *dto.ReviewQueryParams) ([]*reviewmodel.Review, int, error)
	GetReviewByID(id uint) (*reviewmodel.Review, error)
	CreateReview(review *reviewmodel.Review) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(dataSource *db.DataSource) ReviewRepository {
	return &reviewRepository{
		db: dataSource.Db,
	}
}
