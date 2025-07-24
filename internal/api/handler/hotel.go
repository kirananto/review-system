package handler

import (
	"encoding/json"
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
func (h *HotelHandler) CreateHotel(w http.ResponseWriter, r *http.Request) {
	var hotelDto dto.HotelRequestBody
	if err := json.NewDecoder(r.Body).Decode(&hotelDto); err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid request body")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	hotel, errDetails := h.service.CreateHotel(&hotelDto)
	if errDetails != nil {
		errResp := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		response.WriteHTTPResponse(w, errDetails.Code, errResp)
		return
	}

	resp := &response.HTTPResponse{
		Content: hotel,
	}

	response.WriteHTTPResponse(w, http.StatusCreated, resp)
}

func (h *HotelHandler) UpdateHotel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid hotel ID")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	var hotelDto dto.HotelRequestBody
	if err := json.NewDecoder(r.Body).Decode(&hotelDto); err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid request body")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	hotel, errDetails := h.service.UpdateHotel(uint(id), &hotelDto)
	if errDetails != nil {
		errResp := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		response.WriteHTTPResponse(w, errDetails.Code, errResp)
		return
	}

	resp := &response.HTTPResponse{
		Content: hotel,
	}

	response.WriteHTTPResponse(w, http.StatusOK, resp)
}

func (h *HotelHandler) DeleteHotel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid hotel ID")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	errDetails := h.service.DeleteHotel(uint(id))
	if errDetails != nil {
		errResp := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		response.WriteHTTPResponse(w, errDetails.Code, errResp)
		return
	}

	response.WriteHTTPResponse(w, http.StatusNoContent, nil)
}
