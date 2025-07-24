package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/kirananto/review-system/internal/api/handler"
	"github.com/kirananto/review-system/internal/api/service/mock"
	"github.com/kirananto/review-system/internal/logger"
	"github.com/kirananto/review-system/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHotelHandler_GetHotel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService := mock.NewMockHotelService(ctrl)
		log := logger.NewLogger(&logger.LogConfig{LogLevel: "info"})
		hotelHandler := handler.NewHotelHandler(mockService, log)

		expectedHotel := &models.Hotel{
			ID:        1,
			HotelName: "Test Hotel",
		}

		mockService.EXPECT().GetHotelByID(gomock.Any(), uint(1)).Return(expectedHotel, nil)

		req, err := http.NewRequest("GET", "/hotels/1", nil)
		assert.NoError(t, err)

		// We need to set the URL variable for mux to find it
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		rr := httptest.NewRecorder()

		// Act
		hotelHandler.GetHotel(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)

		var actualHotel models.Hotel
		err = json.Unmarshal(rr.Body.Bytes(), &actualHotel)
		assert.NoError(t, err)

		assert.Equal(t, expectedHotel.ID, actualHotel.ID)
		assert.Equal(t, expectedHotel.HotelName, actualHotel.HotelName)
	})

	t.Run("not_found", func(t *testing.T) {
		// Arrange
		mockService := mock.NewMockHotelService(ctrl)
		log := logger.NewLogger(&logger.LogConfig{LogLevel: "info"})
		hotelHandler := handler.NewHotelHandler(mockService, log)

		mockService.EXPECT().GetHotelByID(gomock.Any(), uint(1)).Return(nil, errors.New("not found"))

		req, err := http.NewRequest("GET", "/hotels/1", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		rr := httptest.NewRecorder()

		// Act
		hotelHandler.GetHotel(rr, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestHotelHandler_CreateHotel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService := mock.NewMockHotelService(ctrl)
		log := logger.NewLogger(&logger.LogConfig{LogLevel: "info"})
		hotelHandler := handler.NewHotelHandler(mockService, log)

		newHotel := &models.Hotel{
			HotelName: "Test Hotel",
		}

		mockService.EXPECT().CreateHotel(gomock.Any(), gomock.Any()).Return(nil)

		body, err := json.Marshal(newHotel)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/hotels", bytes.NewReader(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		// Act
		hotelHandler.CreateHotel(rr, req)

		// Assert
		assert.Equal(t, http.StatusCreated, rr.Code)
	})
}

func TestHotelHandler_UpdateHotel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService := mock.NewMockHotelService(ctrl)
		log := logger.NewLogger(&logger.LogConfig{LogLevel: "info"})
		hotelHandler := handler.NewHotelHandler(mockService, log)

		updatedHotel := &models.Hotel{
			ID:        1,
			HotelName: "Updated Test Hotel",
		}

		mockService.EXPECT().UpdateHotel(gomock.Any(), gomock.Any()).Return(nil)

		body, err := json.Marshal(updatedHotel)
		assert.NoError(t, err)

		req, err := http.NewRequest("PUT", "/hotels/1", bytes.NewReader(body))
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		rr := httptest.NewRecorder()

		// Act
		hotelHandler.UpdateHotel(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestHotelHandler_DeleteHotel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService := mock.NewMockHotelService(ctrl)
		log := logger.NewLogger(&logger.LogConfig{LogLevel: "info"})
		hotelHandler := handler.NewHotelHandler(mockService, log)

		mockService.EXPECT().DeleteHotel(gomock.Any(), uint(1)).Return(nil)

		req, err := http.NewRequest("DELETE", "/hotels/1", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		rr := httptest.NewRecorder()

		// Act
		hotelHandler.DeleteHotel(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
