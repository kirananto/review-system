package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer secret",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing auth header",
			authHeader:     "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid auth header format",
			authHeader:     "invalid-format",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "incorrect token",
			authHeader:     "Bearer wrong-secret",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			Auth(handler).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
