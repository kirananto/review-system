package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/kirananto/review-system/internal/api/dto"
	"github.com/kirananto/review-system/internal/api/response"
	"github.com/kirananto/review-system/internal/api/service"
	"github.com/kirananto/review-system/internal/api/utils"
	"github.com/kirananto/review-system/internal/logger"
)

type HotelHandler struct {
	service service.HotelService
	logger  *logger.Logger
	decoder *schema.Decoder
}

func NewHotelHandler(service service.HotelService, logger *logger.Logger) *HotelHandler {
	return &HotelHandler{
		service: service,
		logger:  logger,
		decoder: schema.NewDecoder(),
	}
}

func (h *HotelHandler) GetHotelsList(w http.ResponseWriter, r *http.Request) {
	// Initialize with default values
	queryParams := &dto.HotelsQueryParams{
		Limit:  20,
		Offset: 0,
	}

	// Parse query parameters automatically
	if err := h.decoder.Decode(queryParams, r.URL.Query()); err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid query parameters")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	hotels, total, errorDetails := h.service.GetHotelsList(queryParams)
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
		Results:  hotels,
	}
	resp := &response.HTTPResponse{
		Content: content,
	}

	response.WriteHTTPResponse(w, http.StatusOK, resp)
}

func (h *HotelHandler) GetHotel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid hotel ID")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	hotel, errorDetails := h.service.GetHotelByID(uint(id))
	if errorDetails != nil {
		errResp := response.GetErrorHTTPResponseBody(errorDetails.Code, errorDetails.Message)
		response.WriteHTTPResponse(w, errorDetails.Code, errResp)
		return
	}

	// Create success response
	resp := &response.HTTPResponse{
		Content: hotel,
	}

	response.WriteHTTPResponse(w, http.StatusOK, resp)
}
