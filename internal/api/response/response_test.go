package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteHTTPResponse(t *testing.T) {
	t.Run("valid response", func(t *testing.T) {
		rr := httptest.NewRecorder()
		responseBody := &HTTPResponse{
			Code:    http.StatusOK,
			Message: "Success",
		}
		err := WriteHTTPResponse(rr, http.StatusOK, responseBody)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)

		var actualResponseBody HTTPResponse
		err = json.Unmarshal(rr.Body.Bytes(), &actualResponseBody)
		assert.NoError(t, err)
		assert.Equal(t, responseBody.Code, actualResponseBody.Code)
		assert.Equal(t, responseBody.Message, actualResponseBody.Message)
	})

	t.Run("invalid status code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		err := WriteHTTPResponse(rr, 99, nil)
		assert.Error(t, err)
		err = WriteHTTPResponse(rr, 601, nil)
		assert.Error(t, err)
	})
}

func TestGetErrorHTTPResponseBody(t *testing.T) {
	responseBody := GetErrorHTTPResponseBody(http.StatusBadRequest, "Bad Request")
	assert.Equal(t, http.StatusBadRequest, responseBody.Code)
	assert.Equal(t, "Bad Request", responseBody.Message)
	assert.Equal(t, map[string]interface{}{}, responseBody.Content)
}
