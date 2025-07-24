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
	CreateHotel(hotel *dto.HotelRequestBody) (*models.Hotel, *response.ErrorDetails)
	UpdateHotel(id uint, hotel *dto.HotelRequestBody) (*models.Hotel, *response.ErrorDetails)
	DeleteHotel(id uint) *response.ErrorDetails
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

func (s *hotelService) CreateHotel(hotelDto *dto.HotelRequestBody) (*models.Hotel, *response.ErrorDetails) {
	hotel := &models.Hotel{
		HotelName: hotelDto.HotelName,
	}
	err := s.repo.CreateHotel(hotel)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create hotel",
			Error:   err,
		}
	}
	return hotel, nil
}

func (s *hotelService) UpdateHotel(id uint, hotelDto *dto.HotelRequestBody) (*models.Hotel, *response.ErrorDetails) {
	hotel, errDetails := s.GetHotelByID(id)
	if errDetails != nil {
		return nil, errDetails
	}

	hotel.HotelName = hotelDto.HotelName
	err := s.repo.UpdateHotel(hotel)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update hotel",
			Error:   err,
		}
	}
	return hotel, nil
}

func (s *hotelService) DeleteHotel(id uint) *response.ErrorDetails {
	_, errDetails := s.GetHotelByID(id)
	if errDetails != nil {
		return errDetails
	}

	err := s.repo.DeleteHotel(id)
	if err != nil {
		return &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete hotel",
			Error:   err,
		}
	}
	return nil
}
