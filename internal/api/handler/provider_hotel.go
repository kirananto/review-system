package handler

import (
	"net/http"

	"github.com/gorilla/schema"
	"github.com/kirananto/review-system/internal/api/dto"
	"github.com/kirananto/review-system/internal/api/response"
	"github.com/kirananto/review-system/internal/api/service"
	"github.com/kirananto/review-system/internal/api/utils"
	"github.com/kirananto/review-system/internal/logger"
)

type ProviderHotelHandler struct {
	service service.ProviderHotelService
	logger  *logger.Logger
	decoder *schema.Decoder
}

func NewProviderHotelHandler(service service.ProviderHotelService, logger *logger.Logger) *ProviderHotelHandler {
	return &ProviderHotelHandler{
		service: service,
		logger:  logger,
		decoder: schema.NewDecoder(),
	}
}

func (h *ProviderHotelHandler) GetProviderHotelsList(w http.ResponseWriter, r *http.Request) {
	// Initialize with default values
	queryParams := &dto.ProviderHotelsQueryParams{
		Limit:  20,
		Offset: 0,
	}

	// Parse query parameters automatically
	if err := h.decoder.Decode(queryParams, r.URL.Query()); err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid query parameters")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	providerHotels, total, errorDetails := h.service.GetProviderHotelsList(queryParams)
	if errorDetails != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusInternalServerError, errorDetails.Error.Error())
		response.WriteHTTPResponse(w, http.StatusInternalServerError, errResp)
		return
	}

	// Get pagination links
	prevURL, nextURL := utils.GetPaginationLinks(r, queryParams.Offset, queryParams.Limit, total)

	// Create success response with pagination
	content := &response.HTTPResponseContent{
		Count:    total,
		Previous: prevURL,
		Next:     nextURL,
		Results:  providerHotels,
	}
	resp := &response.HTTPResponse{
		Content: content,
	}

	response.WriteHTTPResponse(w, http.StatusOK, resp)
}
