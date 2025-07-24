package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kirananto/review-system/internal/api/dto"
	"github.com/kirananto/review-system/internal/api/repository"
	"github.com/kirananto/review-system/internal/api/response"
	"github.com/kirananto/review-system/internal/logger"
	"github.com/kirananto/review-system/internal/models"
	"gorm.io/gorm"
)

type ReviewService interface {
	GetReviewsList(queryParam *dto.ReviewQueryParams) ([]*models.Review, int, *response.ErrorDetails)
	GetReviewByID(id uint) (*models.Review, *response.ErrorDetails)
	ProcessReviews(ctx context.Context, reader io.Reader, fileName string) error
}

type reviewService struct {
	repo   repository.ReviewRepository
	logger *logger.Logger
}

func NewReviewService(repo repository.ReviewRepository, logger *logger.Logger) ReviewService {
	return &reviewService{
		repo:   repo,
		logger: logger,
	}
}

func (s *reviewService) GetReviewsList(queryParam *dto.ReviewQueryParams) ([]*models.Review, int, *response.ErrorDetails) {
	reviews, total, err := s.repo.GetReviewsList(queryParam)
	if err != nil {
		return nil, 0, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
		}
	}

	return reviews, total, nil
}

func (s *reviewService) GetReviewByID(id uint) (*models.Review, *response.ErrorDetails) {
	review, err := s.repo.GetReviewByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &response.ErrorDetails{
				Code:    http.StatusNotFound,
				Message: "Review not found",
				Error:   err,
			}
		}
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
		}
	}

	return review, nil
}

type ReviewData struct {
	HotelID   int    `json:"hotelId"`
	Platform  string `json:"platform"`
	HotelName string `json:"hotelName"`
	Comment   struct {
		HotelReviewID      int     `json:"hotelReviewId"`
		Rating             float64 `json:"rating"`
		ReviewComments     string  `json:"reviewComments"`
		ReviewTitle        string  `json:"reviewTitle"`
		ReviewDate         string  `json:"reviewDate"`
		ReviewProviderText string  `json:"reviewProviderText"`
	} `json:"comment"`
	OverallByProviders []struct {
		ProviderID   int     `json:"providerId"`
		Provider     string  `json:"provider"`
		OverallScore float64 `json:"overallScore"`
		ReviewCount  int     `json:"reviewCount"`
		Grades       struct {
			Cleanliness           float64 `json:"Cleanliness"`
			Facilities            float64 `json:"Facilities"`
			Location              float64 `json:"Location"`
			RoomComfortAndQuality float64 `json:"Room comfort and quality"`
			Service               float64 `json:"Service"`
			ValueForMoney         float64 `json:"Value for money"`
		} `json:"grades"`
	} `json:"overallByProviders"`
}

func (s *reviewService) ProcessReviews(ctx context.Context, reader io.Reader, fileName string) error {
	scanner := bufio.NewScanner(reader)
	log := s.logger
	var successCount, failureCount, totalCount int

	for scanner.Scan() {
		totalCount++
		var data ReviewData
		line := scanner.Bytes()

		if err := json.Unmarshal(line, &data); err != nil {
			log.Error(err, fmt.Sprintf("Failed to parse JSON line: %v. Line: %s", err, string(line)))
			failureCount++
			continue
		}

		if err := s.validateData(&data); err != nil {
			log.Error(err, fmt.Sprintf("Invalid data: %v. Data: %+v", err, data))
			failureCount++
			continue
		}

		if err := s.processRecord(ctx, &data); err != nil {
			log.Error(err, fmt.Sprintf("Failed to process record: %v. Data: %+v", err, data))
			failureCount++
			continue
		}
		successCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	auditLog := &models.AuditLog{
		FileName:     fileName,
		SuccessCount: successCount,
		FailureCount: failureCount,
		TotalCount:   totalCount,
	}

	if err := s.repo.CreateAuditLog(auditLog); err != nil {
		log.Error(err, "Failed to create audit log")
		// Do not return error, as the main process was successful
	}

	log.Info(fmt.Sprintf("Processed file: %s, Success: %d, Failed: %d, Total: %d", fileName, successCount, failureCount, totalCount))

	return nil
}

func (s *reviewService) validateData(data *ReviewData) error {
	if data.Comment.HotelReviewID == 0 {
		return fmt.Errorf("HotelReviewID is required")
	}
	if data.Platform == "" {
		return fmt.Errorf("platform is required")
	}
	if data.HotelName == "" {
		return fmt.Errorf("hotelName is required")
	}
	if data.HotelID == 0 {
		return fmt.Errorf("hotelId is required")
	}
	return nil
}

func (s *reviewService) processRecord(ctx context.Context, data *ReviewData) error {
	log := s.logger
	// Get or create provider
	provider, err := s.getOrCreateProvider(data.Comment.ReviewProviderText)
	if err != nil {
		return err
	}

	// Get or create hotel
	hotel, err := s.getOrCreateHotel(data.HotelName)
	if err != nil {
		return err
	}

	// Find the correct provider data from the OverallByProviders array
	var providerData struct {
		OverallScore float64     `json:"overallScore"`
		ReviewCount  int         `json:"reviewCount"`
		Grades       interface{} `json:"grades"`
	}
	for _, p := range data.OverallByProviders {
		if p.Provider == data.Platform {
			providerData.OverallScore = p.OverallScore
			providerData.ReviewCount = p.ReviewCount
			providerData.Grades = p.Grades
			break
		}
	}

	gradesJSON, err := json.Marshal(providerData.Grades)
	if err != nil {
		return fmt.Errorf("failed to marshal grades: %w", err)
	}

	// Get or create provider-hotel mapping
	if _, err := s.getOrCreateProviderHotel(provider.ID, hotel.ID, providerData.OverallScore, providerData.ReviewCount, string(gradesJSON)); err != nil {
		return err
	}

	reviewDate, err := time.Parse(time.RFC3339, data.Comment.ReviewDate)
	if err != nil {
		log.Info(fmt.Sprintf("Could not parse review date: %v", err))
		reviewDate = time.Now()
	}

	// Create review
	review := &models.Review{
		ProviderID:   provider.ID,
		HotelID:      hotel.ID,
		ID:           uint(data.Comment.HotelReviewID),
		Rating:       data.Comment.Rating,
		Comment:      data.Comment.ReviewComments,
		Lang:         "en",
		ReviewDate:   reviewDate,
		ReviewerInfo: []byte(`{}`), // Empty JSON object since we don't have reviewer info
	}

	if err := s.repo.UpsertReview(review); err != nil {
		return fmt.Errorf("failed to create or update review: %w", err)
	}

	return nil
}

func (s *reviewService) getOrCreateProvider(name string) (*models.Provider, error) {
	provider, err := s.repo.GetProviderByName(name)
	if err == nil {
		return provider, nil
	}

	provider = &models.Provider{Name: name}
	if err := s.repo.CreateProvider(provider); err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	return provider, nil
}

func (s *reviewService) getOrCreateHotel(name string) (*models.Hotel, error) {
	hotel, err := s.repo.GetHotelByName(name)
	if err == nil {
		return hotel, nil
	}

	hotel = &models.Hotel{HotelName: name}
	if err := s.repo.CreateHotel(hotel); err != nil {
		return nil, fmt.Errorf("failed to create hotel: %w", err)
	}

	return hotel, nil
}

func (s *reviewService) getOrCreateProviderHotel(providerID, hotelID uint, overallScore float64, reviewCount int, gradesJSON string) (*models.ProviderHotel, error) {
	providerHotel, err := s.repo.GetProviderHotel(providerID, hotelID)
	if err == nil {
		// Update existing record
		providerHotel.OverallScore = overallScore
		providerHotel.ReviewCount = reviewCount
		providerHotel.Grades = []byte(gradesJSON)
		if err := s.repo.UpdateProviderHotel(providerHotel); err != nil {
			return nil, fmt.Errorf("failed to update provider hotel: %w", err)
		}
		return providerHotel, nil
	}

	// Create new record
	providerHotel = &models.ProviderHotel{
		ProviderID:   providerID,
		HotelID:      hotelID,
		OverallScore: overallScore,
		ReviewCount:  reviewCount,
		Grades:       []byte(gradesJSON),
	}

	if err := s.repo.CreateProviderHotel(providerHotel); err != nil {
		return nil, fmt.Errorf("failed to create provider hotel: %w", err)
	}

	return providerHotel, nil
}
