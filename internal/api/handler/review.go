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

// ReviewHandler handles API requests for reviews
type ReviewHandler struct {
	service service.ReviewService
	logger  *logger.Logger
	decoder *schema.Decoder
}

// NewReviewHandler creates a new ReviewHandler
func NewReviewHandler(service service.ReviewService, logger *logger.Logger) *ReviewHandler {
	return &ReviewHandler{
		service: service,
		logger:  logger,
		decoder: schema.NewDecoder(),
	}
}

func (h *ReviewHandler) GetReviewsList(w http.ResponseWriter, r *http.Request) {
	// Initialize with default values
	queryParams := &dto.ReviewQueryParams{
		Limit:  20,
		Offset: 0,
	}

	// Parse query parameters automatically
	if err := h.decoder.Decode(queryParams, r.URL.Query()); err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid query parameters")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	reviews, total, errorDetails := h.service.GetReviewsList(queryParams)
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
		Results:  reviews,
	}
	resp := &response.HTTPResponse{
		Content: content,
	}

	response.WriteHTTPResponse(w, http.StatusOK, resp)
}

// GetReview godoc
// @Summary Get a review by ID
// @Description Get a review by ID
// @ID get-review-by-id
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} review.Review
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}
	review, err := h.service.GetReviewByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// CreateReview godoc
// @Summary Create a new review
// @Description Create a new review
// @ID create-review
// @Accept json
// @Param review body review.Review true "Review object"
// @Success 201
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var review review.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateReview(r.Context(), &review); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// UpdateReview godoc
// @Summary Update an existing review
// @Description Update an existing review
// @ID update-review
// @Accept json
// @Param id path int true "Review ID"
// @Param review body review.Review true "Review object"
// @Success 200
// @Router /reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}
	var review review.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	review.ID = uint(id)
	if err := h.service.UpdateReview(r.Context(), &review); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteReview godoc
// @Summary Delete a review
// @Description Delete a review
// @ID delete-review
// @Param id path int true "Review ID"
// @Success 200
// @Router /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteReview(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
