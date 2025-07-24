package service

import (
	"net/http"

	"github.com/kirananto/review-system/internal/api/dto"
	"github.com/kirananto/review-system/internal/api/repository"
	"github.com/kirananto/review-system/internal/api/response"
	"github.com/kirananto/review-system/internal/logger"
	"github.com/kirananto/review-system/pkg/review"
)

type ProviderHotelService interface {
	GetProviderHotelsList(queryParam *dto.ProviderHotelsQueryParams) ([]*review.ProviderHotel, int, *response.ErrorDetails)
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
