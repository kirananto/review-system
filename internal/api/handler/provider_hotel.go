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
	"github.com/kirananto/review-system/pkg/review"
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

// GetProviderHotel godoc
// @Summary Get a provider hotel by ID
// @Description Get a provider hotel by ID
// @ID get-provider-hotel-by-id
// @Produce json
// @Param id path int true "Provider Hotel ID"
// @Success 200 {object} review.ProviderHotel
// @Router /provider-hotels/{id} [get]
func (h *ProviderHotelHandler) GetProviderHotel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid provider hotel ID", http.StatusBadRequest)
		return
	}
	providerHotel, err := h.service.GetProviderHotelByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providerHotel)
}

// CreateProviderHotel godoc
// @Summary Create a new provider hotel
// @Description Create a new provider hotel
// @ID create-provider-hotel
// @Accept json
// @Param providerHotel body review.ProviderHotel true "Provider Hotel object"
// @Success 201
// @Router /provider-hotels [post]
func (h *ProviderHotelHandler) CreateProviderHotel(w http.ResponseWriter, r *http.Request) {
	var providerHotel review.ProviderHotel
	if err := json.NewDecoder(r.Body).Decode(&providerHotel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateProviderHotel(r.Context(), &providerHotel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// UpdateProviderHotel godoc
// @Summary Update an existing provider hotel
// @Description Update an existing provider hotel
// @ID update-provider-hotel
// @Accept json
// @Param id path int true "Provider Hotel ID"
// @Param providerHotel body review.ProviderHotel true "Provider Hotel object"
// @Success 200
// @Router /provider-hotels/{id} [put]
func (h *ProviderHotelHandler) UpdateProviderHotel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid provider hotel ID", http.StatusBadRequest)
		return
	}
	var providerHotel review.ProviderHotel
	if err := json.NewDecoder(r.Body).Decode(&providerHotel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	providerHotel.ID = uint(id)
	if err := h.service.UpdateProviderHotel(r.Context(), &providerHotel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteProviderHotel godoc
// @Summary Delete a provider hotel
// @Description Delete a provider hotel
// @ID delete-provider-hotel
// @Param id path int true "Provider Hotel ID"
// @Success 200
// @Router /provider-hotels/{id} [delete]
func (h *ProviderHotelHandler) DeleteProviderHotel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid provider hotel ID", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteProviderHotel(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
