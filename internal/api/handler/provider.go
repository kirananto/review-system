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

type ProviderHandler struct {
	service service.ProviderService
	logger  *logger.Logger
	decoder *schema.Decoder
}

func NewProviderHandler(service service.ProviderService, logger *logger.Logger) *ProviderHandler {
	return &ProviderHandler{
		service: service,
		logger:  logger,
		decoder: schema.NewDecoder(),
	}
}

// GetProvidersList godoc
// @Summary Get a list of providers
// @Description Get a list of providers with optional filters
// @ID get-providers-list
// @Produce json
// @Param name query string false "Provider name"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.HTTPResponse{content=response.HTTPResponseContent{results=[]models.Provider}}
// @Router /providers [get]
func (h *ProviderHandler) GetProvidersList(w http.ResponseWriter, r *http.Request) {
	// Initialize with default values
	queryParams := &dto.ProvidersQueryParams{
		Limit:  20,
		Offset: 0,
	}

	// Parse query parameters automatically
	if err := h.decoder.Decode(queryParams, r.URL.Query()); err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid query parameters")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	providers, total, errorDetails := h.service.GetProvidersList(queryParams)
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
		Results:  providers,
	}
	resp := &response.HTTPResponse{
		Content: content,
	}

	response.WriteHTTPResponse(w, http.StatusOK, resp)
}

// GetProvider godoc
// @Summary Get a provider by ID
// @Description Get a provider by ID
// @ID get-provider-by-id
// @Produce json
// @Param id path int true "Provider ID"
// @Success 200 {object} response.HTTPResponse{content=models.Provider}
// @Router /providers/{id} [get]
func (h *ProviderHandler) GetProvider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		errResp := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid Provider ID")
		response.WriteHTTPResponse(w, http.StatusBadRequest, errResp)
		return
	}

	provider, errorDetails := h.service.GetProviderByID(uint(id))
	if errorDetails != nil {
		errResp := response.GetErrorHTTPResponseBody(errorDetails.Code, errorDetails.Message)
		response.WriteHTTPResponse(w, errorDetails.Code, errResp)
		return
	}

	// Create success response
	resp := &response.HTTPResponse{
		Content: provider,
	}

	response.WriteHTTPResponse(w, http.StatusOK, resp)
}
