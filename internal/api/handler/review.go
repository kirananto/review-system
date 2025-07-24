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

// GetReviewsList godoc
// @Summary Get a list of reviews
// @Description Get a list of reviews with optional filters
// @ID get-reviews-list
// @Produce json
// @Param hotel_id query int false "Hotel ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.HTTPResponse{content=response.HTTPResponseContent{results=[]review.Review}}
// @Router /reviews [get]
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
// @Summary Get a Review by ID
// @Description Get a Review by ID
// @ID get-review-by-id
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} review.Review
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid Review ID")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	review, errorDetails := h.service.GetReviewByID(uint(id))
	if errorDetails != nil {
		errResp := response.GetErrorHTTPResponseBody(errorDetails.Code, errorDetails.Message)
		response.WriteHTTPResponse(w, errorDetails.Code, errResp)
		return
	}

	// Create success response
	resp := &response.HTTPResponse{
		Content: review,
	}

	response.WriteHTTPResponse(w, http.StatusOK, resp)
}
