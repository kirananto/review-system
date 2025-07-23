package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPResponse struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

type HTTPResponseContent struct {
	Count    int         `json:"count"`
	Previous *string     `json:"prev"`
	Next     *string     `json:"next"`
	Results  interface{} `json:"results"`
}

// Error details will be returned by a service function to the handler
type ErrorDetails struct {
	Code    int
	Message string
	Error   error
}

// StatusCode refer to Http Code
func WriteHTTPResponse(w http.ResponseWriter, statusCode int, responseBody *HTTPResponse) error {
	if statusCode < 100 || statusCode > 600 {
		return fmt.Errorf("invalid status code for HTTP response: %v", statusCode)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(responseBody)
}

// Code refer to Application Code
func GetErrorHTTPResponseBody(code int, message string) *HTTPResponse {
	return &HTTPResponse{
		Code:    code,
		Message: message,
		Content: map[string]interface{}{},
	}
}
