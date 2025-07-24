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
	"github.com/kirananto/review-system/pkg/review"
)

type ReviewService interface {
	GetReviewsList(queryParam *dto.ReviewQueryParams) ([]*review.Review, int, *response.ErrorDetails)
	GetReviewByID(ctx context.Context, id uint) (*review.Review, error)
	CreateReview(ctx context.Context, review *review.Review) error
	UpdateReview(ctx context.Context, review *review.Review) error
	DeleteReview(ctx context.Context, id uint) error
	ProcessReviews(ctx context.Context, reader io.Reader) error
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

func (s *reviewService) GetReviewsList(queryParam *dto.ReviewQueryParams) ([]*review.Review, int, *response.ErrorDetails) {
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

func (s *reviewService) GetReviewByID(ctx context.Context, id uint) (*review.Review, error) {
	return s.repo.GetReviewByID(id)
}

func (s *reviewService) CreateReview(ctx context.Context, review *review.Review) error {
	return s.repo.CreateReview(review)
}

func (s *reviewService) UpdateReview(ctx context.Context, review *review.Review) error {
	return s.repo.UpdateReview(review)
}

func (s *reviewService) DeleteReview(ctx context.Context, id uint) error {
	return s.repo.DeleteReview(id)
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

func (s *reviewService) ProcessReviews(ctx context.Context, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	log := s.logger
	for scanner.Scan() {
		var data ReviewData
		line := scanner.Bytes()

		if err := json.Unmarshal(line, &data); err != nil {
			log.Error(err, fmt.Sprintf("Failed to parse JSON line: %v. Line: %s", err, string(line)))
			continue
		}

		if err := s.validateData(&data); err != nil {
			log.Error(err, fmt.Sprintf("Invalid data: %v. Data: %+v", err, data))
			continue
		}

		if err := s.processRecord(ctx, &data); err != nil {
			log.Error(err, fmt.Sprintf("Failed to process record: %v. Data: %+v", err, data))
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

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
	providerHotel, err := s.getOrCreateProviderHotel(provider.ID, hotel.ID, fmt.Sprintf("%d", data.HotelID), providerData.OverallScore, providerData.ReviewCount, string(gradesJSON))
	if err != nil {
		return err
	}

	reviewDate, err := time.Parse(time.RFC3339, data.Comment.ReviewDate)
	if err != nil {
		log.Info(fmt.Sprintf("Could not parse review date: %v", err))
		reviewDate = time.Now()
	}

	// Create review
	review := &review.Review{
		ProviderID:       provider.ID,
		ProviderReviewID: fmt.Sprintf("%d", data.Comment.HotelReviewID),
		ProviderHotelID:  providerHotel.ID,
		Rating:           data.Comment.Rating,
		Comment:          data.Comment.ReviewComments,
		Lang:             "en",
		ReviewDate:       reviewDate,
		ReviewerInfo:     "{}", // Empty JSON object since we don't have reviewer info
	}

	if err := s.repo.CreateReview(review); err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

func (s *reviewService) getOrCreateProvider(name string) (*review.Provider, error) {
	provider, err := s.repo.GetProviderByName(name)
	if err == nil {
		return provider, nil
	}

	provider = &review.Provider{Name: name}
	if err := s.repo.CreateProvider(provider); err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	return provider, nil
}

func (s *reviewService) getOrCreateHotel(name string) (*review.Hotel, error) {
	hotel, err := s.repo.GetHotelByName(name)
	if err == nil {
		return hotel, nil
	}

	hotel = &review.Hotel{HotelName: name}
	if err := s.repo.CreateHotel(hotel); err != nil {
		return nil, fmt.Errorf("failed to create hotel: %w", err)
	}

	return hotel, nil
}

func (s *reviewService) getOrCreateProviderHotel(providerID, hotelID uint, providerHotelID string, overallScore float64, reviewCount int, gradesJSON string) (*review.ProviderHotel, error) {
	providerHotel, err := s.repo.GetProviderHotel(providerID, providerHotelID)
	if err == nil {
		// Update existing record
		providerHotel.OverallScore = overallScore
		providerHotel.ReviewCount = reviewCount
		providerHotel.Grades = gradesJSON
		if err := s.repo.UpdateProviderHotel(providerHotel); err != nil {
			return nil, fmt.Errorf("failed to update provider hotel: %w", err)
		}
		return providerHotel, nil
	}

	// Create new record
	providerHotel = &review.ProviderHotel{
		ProviderID:      providerID,
		HotelID:         hotelID,
		ProviderHotelID: providerHotelID,
		OverallScore:    overallScore,
		ReviewCount:     reviewCount,
		Grades:          gradesJSON,
	}

	if err := s.repo.CreateProviderHotel(providerHotel); err != nil {
		return nil, fmt.Errorf("failed to create provider hotel: %w", err)
	}

	return providerHotel, nil
}
