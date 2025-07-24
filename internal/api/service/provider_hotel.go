package service

import (
	"context"
	"net/http"

	"github.com/kirananto/review-system/internal/api/dto"
	"github.com/kirananto/review-system/internal/api/repository"
	"github.com/kirananto/review-system/internal/api/response"
	"github.com/kirananto/review-system/internal/logger"
	"github.com/kirananto/review-system/pkg/review"
)

type ProviderHotelService interface {
	GetProviderHotelsList(queryParam *dto.ProviderHotelsQueryParams) ([]*review.ProviderHotel, int, *response.ErrorDetails)
	GetProviderHotelByID(ctx context.Context, id uint) (*review.ProviderHotel, error)
	CreateProviderHotel(ctx context.Context, providerHotel *review.ProviderHotel) error
	UpdateProviderHotel(ctx context.Context, providerHotel *review.ProviderHotel) error
	DeleteProviderHotel(ctx context.Context, id uint) error
}

type providerHotelService struct {
	repo   repository.ReviewRepository
	logger *logger.Logger
}

func NewProviderHotelService(repo repository.ReviewRepository, logger *logger.Logger) ProviderHotelService {
	return &providerHotelService{
		repo:   repo,
		logger: logger,
	}
}

func (s *providerHotelService) GetProviderHotelsList(queryParam *dto.ProviderHotelsQueryParams) ([]*review.ProviderHotel, int, *response.ErrorDetails) {
	providerHotels, total, err := s.repo.GetProviderHotelsList(queryParam)
	if err != nil {
		return nil, 0, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
		}
	}

	return providerHotels, total, nil
}

func (s *providerHotelService) GetProviderHotelByID(ctx context.Context, id uint) (*review.ProviderHotel, error) {
	return s.repo.GetProviderHotelByID(id)
}

func (s *providerHotelService) CreateProviderHotel(ctx context.Context, providerHotel *review.ProviderHotel) error {
	return s.repo.CreateProviderHotel(providerHotel)
}

func (s *providerHotelService) UpdateProviderHotel(ctx context.Context, providerHotel *review.ProviderHotel) error {
	return s.repo.UpdateProviderHotel(providerHotel)
}

func (s *providerHotelService) DeleteProviderHotel(ctx context.Context, id uint) error {
	return s.repo.DeleteProviderHotel(id)
}
