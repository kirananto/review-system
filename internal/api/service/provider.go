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

type ProviderService interface {
	GetProvidersList(queryParams *dto.ProvidersQueryParams) ([]*models.Provider, int, *response.ErrorDetails)
	GetProviderByID(id uint) (*models.Provider, *response.ErrorDetails)
}

type providerService struct {
	repo   repository.ReviewRepository
	logger *logger.Logger
}

func NewProviderService(repo repository.ReviewRepository, logger *logger.Logger) ProviderService {
	return &providerService{
		repo:   repo,
		logger: logger,
	}
}

func (s *providerService) GetProvidersList(queryParams *dto.ProvidersQueryParams) ([]*models.Provider, int, *response.ErrorDetails) {
	providers, total, err := s.repo.GetProvidersList(queryParams)
	if err != nil {
		return nil, 0, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
		}
	}

	return providers, total, nil
}

func (s *providerService) GetProviderByID(id uint) (*models.Provider, *response.ErrorDetails) {
	provider, err := s.repo.GetProviderByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &response.ErrorDetails{
				Code:    http.StatusNotFound,
				Message: "Provider not found",
				Error:   err,
			}
		}
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
		}
	}
	return provider, nil
}
