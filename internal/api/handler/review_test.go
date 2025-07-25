package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/kirananto/review-system/internal/api/handler"
	"github.com/kirananto/review-system/internal/api/response"
	"github.com/kirananto/review-system/internal/api/service/mock"
	"github.com/kirananto/review-system/internal/logger"
	"github.com/kirananto/review-system/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestReviewHandler_GetReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService := mock.NewMockReviewService(ctrl)
		log := logger.NewLogger(&logger.LogConfig{LogLevel: "info"})
		reviewHandler := handler.NewReviewHandler(mockService, log)

		expectedReview := &models.Review{
			ID:         1,
			Comment:    "Great hotel!",
			Rating:     5,
			ReviewDate: time.Now(),
		}

		mockService.EXPECT().GetReviewByID(uint(1)).Return(expectedReview, nil)

		req, err := http.NewRequest("GET", "/reviews/1", nil)
		assert.NoError(t, err)

		// We need to set the URL variable for mux to find it
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		rr := httptest.NewRecorder()

		// Act
		reviewHandler.GetReview(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)

		var resp response.HTTPResponse
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)

		var actualReview models.Review
		actualReviewBytes, _ := json.Marshal(resp.Content)
		err = json.Unmarshal(actualReviewBytes, &actualReview)
		assert.NoError(t, err)

		// We need to ignore the time difference in the assertion
		assert.Equal(t, expectedReview.ID, actualReview.ID)
		assert.Equal(t, expectedReview.Comment, actualReview.Comment)
		assert.Equal(t, expectedReview.Rating, actualReview.Rating)

	})

	t.Run("not_found", func(t *testing.T) {
		// Arrange
		mockService := mock.NewMockReviewService(ctrl)
		log := logger.NewLogger(&logger.LogConfig{LogLevel: "info"})
		reviewHandler := handler.NewReviewHandler(mockService, log)

		mockService.EXPECT().GetReviewByID(uint(1)).Return(nil, &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Review not found",
			Error:   errors.New("not found"),
		})

		req, err := http.NewRequest("GET", "/reviews/1", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		rr := httptest.NewRecorder()

		// Act
		reviewHandler.GetReview(rr, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
