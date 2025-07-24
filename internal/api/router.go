package api

import (
	"github.com/gorilla/mux"
	"github.com/kirananto/review-system/internal/api/handler"
	"github.com/kirananto/review-system/internal/api/repository"
	"github.com/kirananto/review-system/internal/api/service"
	"github.com/kirananto/review-system/internal/db"
	"github.com/kirananto/review-system/internal/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

func getProviderHandler(dataSource *db.DataSource, log *logger.Logger) *handler.ProviderHandler {
	repository := repository.NewReviewRepository(dataSource)
	service := service.NewProviderService(repository, log)
	return handler.NewProviderHandler(service, log)
}

func getHotelHandler(dataSource *db.DataSource, log *logger.Logger) *handler.HotelHandler {
	repository := repository.NewReviewRepository(dataSource)
	service := service.NewHotelService(repository, log)
	return handler.NewHotelHandler(service, log)
}

func getProviderHotelHandler(dataSource *db.DataSource, log *logger.Logger) *handler.ProviderHotelHandler {
	repository := repository.NewReviewRepository(dataSource)
	service := service.NewProviderHotelService(repository, log)
	return handler.NewProviderHotelHandler(service, log)
}

func getReviewHandler(dataSource *db.DataSource, log *logger.Logger) *handler.ReviewHandler {
	repository := repository.NewReviewRepository(dataSource)
	service := service.NewReviewService(repository, log)
	return handler.NewReviewHandler(service, log)
}

func SetUpRoutes(dataSource *db.DataSource, log *logger.Logger) *mux.Router {
	r := mux.NewRouter()

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Initialize handlers
	providerHandler := getProviderHandler(dataSource, log)
	hotelHandler := getHotelHandler(dataSource, log)
	providerHotelHandler := getProviderHotelHandler(dataSource, log)
	reviewHandler := getReviewHandler(dataSource, log)

	// Provider routes
	api.HandleFunc("/providers", providerHandler.GetProvidersList).Methods("GET")
	api.HandleFunc("/providers/{id:[0-9]+}", providerHandler.GetProvider).Methods("GET")

	// Hotel routes
	api.HandleFunc("/hotels", hotelHandler.GetHotelsList).Methods("GET")
	api.HandleFunc("/hotels/{id:[0-9]+}", hotelHandler.GetHotel).Methods("GET")

	// ProviderHotel routes
	api.HandleFunc("/provider-hotels", providerHotelHandler.GetProviderHotelsList).Methods("GET")

	// Review routes
	api.HandleFunc("/reviews", reviewHandler.GetReviews).Methods("GET")
	api.HandleFunc("/reviews/{id:[0-9]+}", reviewHandler.GetReview).Methods("GET")
	api.HandleFunc("/reviews", reviewHandler.CreateReview).Methods("POST")
	api.HandleFunc("/reviews/{id:[0-9]+}", reviewHandler.UpdateReview).Methods("PUT")
	api.HandleFunc("/reviews/{id:[0-9]+}", reviewHandler.DeleteReview).Methods("DELETE")

	return r
}
