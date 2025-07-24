package service

import (
	"net/http"

	"github.com/kirananto/review-system/internal/api/dto"
	"github.com/kirananto/review-system/internal/api/repository"
	"github.com/kirananto/review-system/internal/api/response"
	"github.com/kirananto/review-system/internal/logger"
	"github.com/kirananto/review-system/internal/models"
	"gorm.io/gorm"
)

type HotelService interface {
	GetHotelsList(queryParam *dto.HotelsQueryParams) ([]*models.Hotel, int, *response.ErrorDetails)
	GetHotelByID(id uint) (*models.Hotel, *response.ErrorDetails)
}

type hotelService struct {
	repo   repository.ReviewRepository
	logger *logger.Logger
}

func NewHotelService(repo repository.ReviewRepository, logger *logger.Logger) HotelService {
	return &hotelService{
		repo:   repo,
		logger: logger,
	}
}

func (s *hotelService) GetHotelsList(queryParam *dto.HotelsQueryParams) ([]*models.Hotel, int, *response.ErrorDetails) {
	hotels, total, err := s.repo.GetHotelsList(queryParam)
	if err != nil {
		return nil, 0, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
		}
	}

	return hotels, total, nil
}

func (s *hotelService) GetHotelByID(id uint) (*models.Hotel, *response.ErrorDetails) {
	hotels, err := s.repo.GetHotelByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &response.ErrorDetails{
				Code:    http.StatusNotFound,
				Message: "Hotel not found",
				Error:   err,
			}
		}
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
		}
	}

	return hotels, nil
}
