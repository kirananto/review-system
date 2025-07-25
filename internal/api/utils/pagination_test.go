package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPaginationLinks(t *testing.T) {
	tests := []struct {
		name         string
		offset       int
		limit        int
		total        int
		expectedPrev *string
		expectedNext *string
	}{
		{
			name:         "first page",
			offset:       0,
			limit:        10,
			total:        30,
			expectedPrev: nil,
			expectedNext: stringPtr("/?limit=10&offset=10"),
		},
		{
			name:         "middle page",
			offset:       10,
			limit:        10,
			total:        30,
			expectedPrev: stringPtr("/?limit=10&offset=0"),
			expectedNext: stringPtr("/?limit=10&offset=20"),
		},
		{
			name:         "last page",
			offset:       20,
			limit:        10,
			total:        30,
			expectedPrev: stringPtr("/?limit=10&offset=10"),
			expectedNext: nil,
		},
		{
			name:         "single page",
			offset:       0,
			limit:        10,
			total:        5,
			expectedPrev: nil,
			expectedNext: nil,
		},
		{
			name:         "prev offset less than 0",
			offset:       5,
			limit:        10,
			total:        20,
			expectedPrev: stringPtr("/?limit=10&offset=0"),
			expectedNext: stringPtr("/?limit=10&offset=15"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			prev, next := GetPaginationLinks(req, tt.offset, tt.limit, tt.total)
			assert.Equal(t, tt.expectedPrev, prev)
			assert.Equal(t, tt.expectedNext, next)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
