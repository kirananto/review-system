package validator

import (
	"testing"

	"github.com/kirananto/review-system/internal/api/dto"
	"github.com/stretchr/testify/assert"
)

func TestReviewValidator_ValidateCreateReview(t *testing.T) {
	validator := NewReviewValidator()

	tests := []struct {
		name        string
		req         *dto.ReviewRequest
		expectedErr string
	}{
		{
			name: "valid request",
			req: &dto.ReviewRequest{
				ProductID:  "123",
				Rating:     4.5,
				ExternalID: "ext123",
			},
			expectedErr: "",
		},
		{
			name: "missing product_id",
			req: &dto.ReviewRequest{
				Rating:     4.5,
				ExternalID: "ext123",
			},
			expectedErr: "product_id is required",
		},
		{
			name: "rating less than 0",
			req: &dto.ReviewRequest{
				ProductID:  "123",
				Rating:     -1,
				ExternalID: "ext123",
			},
			expectedErr: "rating must be between 0 and 5",
		},
		{
			name: "rating greater than 5",
			req: &dto.ReviewRequest{
				ProductID:  "123",
				Rating:     6,
				ExternalID: "ext123",
			},
			expectedErr: "rating must be between 0 and 5",
		},
		{
			name: "missing external_id",
			req: &dto.ReviewRequest{
				ProductID: "123",
				Rating:    4.5,
			},
			expectedErr: "external_id required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateReview(tt.req)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}
