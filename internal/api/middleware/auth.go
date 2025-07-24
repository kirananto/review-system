package middleware

import (
	"net/http"
	"strings"

	"github.com/kirananto/review-system/internal/api/response"
)

// Auth is a middleware that checks for a valid API key.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errResp := response.GetErrorHTTPResponseBody(http.StatusUnauthorized, "Authorization header is required hint: use `Bearer secret`")
			response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errResp := response.GetErrorHTTPResponseBody(http.StatusUnauthorized, "Invalid Authorization header format hint: use `Bearer secret`")
			response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
			return
		}

		token := parts[1]
		if token != "secret" {

			errResp := response.GetErrorHTTPResponseBody(http.StatusUnauthorized, "IInvalid API key hint: use `Bearer secret` ")
			response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
			return
		}

		next.ServeHTTP(w, r)
	})
}
