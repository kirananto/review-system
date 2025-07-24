package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/mux"
	"github.com/kirananto/review-system/internal/api"
	"github.com/kirananto/review-system/internal/api/repository"
	"github.com/kirananto/review-system/internal/api/service"
	"github.com/kirananto/review-system/internal/db"
	"github.com/kirananto/review-system/internal/logger"
	models "github.com/kirananto/review-system/internal/models"
	"github.com/kirananto/review-system/internal/s3"
)

type Server struct {
	Config     *ServerConfig
	Logger     *logger.Logger
	DataSource *db.DataSource
	Router     *mux.Router
	S3Service  s3.S3Service
}
type ServerConfig struct {
	DatabaseDSN string
	RunMode     string // "local" or "lambda"
	Port        string // e.g., ":8000"
	LogConfig   logger.LogConfig
}

// ResponseWriter captures the response for Lambda
type ResponseWriter struct {
	Headers    map[string]string
	Body       strings.Builder
	StatusCode int
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		Headers:    make(map[string]string),
		StatusCode: http.StatusOK,
	}
}

func (w *ResponseWriter) Header() http.Header {
	header := http.Header{}
	for k, v := range w.Headers {
		header.Set(k, v)
	}
	return header
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.Body.Write(b)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	log := logger.NewLogger(&cfg.LogConfig)

	dataSource := db.NewDataSource(cfg.DatabaseDSN)

	dataSource.Db.AutoMigrate(&models.Provider{}, &models.Hotel{}, &models.Review{}, &models.ProviderHotel{})

	router := api.SetUpRoutes(dataSource, log)

	// Initialize S3 service
	s3Service, err := s3.NewS3Service()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 service: %w", err)
	}

	server := &Server{
		Config:     cfg,
		Logger:     log,
		DataSource: dataSource,
		S3Service:  s3Service,
		Router:     router,
	}

	return server, nil
}

func (s *Server) Start() error {
	if s.Config.RunMode == "local" {
		s.Logger.Info(fmt.Sprintf("Starting local server on %s\n", s.Config.Port))
		return http.ListenAndServe(s.Config.Port, s.Router)
	}

	s.Logger.Info("Starting Lambda handler\n")
	lambda.Start(s.handle)
	return nil
}

func (s *Server) handle(ctx context.Context, event json.RawMessage) (interface{}, error) {
	// First, try to unmarshal as an API Gateway request
	var apiReq events.APIGatewayProxyRequest
	if err := json.Unmarshal(event, &apiReq); err == nil && apiReq.RequestContext.RequestID != "" {
		return s.handleAPIGatewayRequest(ctx, apiReq)
	}

	// If that fails, try to unmarshal as an SQS event
	var sqsEvent events.SQSEvent
	if err := json.Unmarshal(event, &sqsEvent); err == nil && len(sqsEvent.Records) > 0 {
		return s.handleSQSEvent(ctx, sqsEvent)
	}

	return nil, fmt.Errorf("unsupported event type")
}

// handleAPIGatewayRequest handles API Gateway proxy requests.
func (s *Server) handleAPIGatewayRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	httpReq, err := s.convertToHTTPRequest(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error converting request: %v", err),
		}, nil
	}

	w := NewResponseWriter()
	s.Router.ServeHTTP(w, httpReq)

	return events.APIGatewayProxyResponse{
		StatusCode: w.StatusCode,
		Headers:    w.Headers,
		Body:       w.Body.String(),
	}, nil
}

func (s *Server) handleSQSEvent(ctx context.Context, sqsEvent events.SQSEvent) (interface{}, error) {
	log := s.Logger
	for _, record := range sqsEvent.Records {
		log.Info(fmt.Sprintf("Processing SQS message: %s", record.MessageId))

		var s3Event events.S3Event
		if err := json.Unmarshal([]byte(record.Body), &s3Event); err != nil {
			log.Info(fmt.Sprintf("Error unmarshalling S3 event from SQS message %s: %v", record.MessageId, err))
			continue
		}

		for _, s3Record := range s3Event.Records {
			bucket := s3Record.S3.Bucket.Name
			key := s3Record.S3.Object.Key

			log.Info(fmt.Sprintf("Processing S3 object: bucket=%s, key=%s", bucket, key))

			reader, err := s.S3Service.GetObject(ctx, bucket, key)
			if err != nil {
				log.Error(err, fmt.Sprintf("Error getting S3 object %s/%s: %v", bucket, key, err))
				continue
			}
			defer reader.Close()

			repository := repository.NewReviewRepository(s.DataSource)
			reviewService := service.NewReviewService(repository, log)

			if err := reviewService.ProcessReviews(ctx, reader, key); err != nil {
				log.Error(err, fmt.Sprintf("Error processing reviews from S3 object %s/%s: %v", bucket, key, err))
				continue
			}
		}
	}

	return nil, nil
}

// Helper to convert API Gateway request to http.Request
func (s *Server) convertToHTTPRequest(req events.APIGatewayProxyRequest) (*http.Request, error) {
	// Use the path directly from the request
	urlPath := req.Path
	if req.QueryStringParameters != nil {
		queryParams := url.Values{}
		for key, value := range req.QueryStringParameters {
			queryParams.Set(key, value)
		}
		urlPath = fmt.Sprintf("%s?%s", urlPath, queryParams.Encode())
	}

	// Create the request body
	var body io.Reader
	if req.Body != "" {
		body = strings.NewReader(req.Body)
	}

	// Create the http.Request
	// In Lambda, the host is not present in the request, so we add a dummy one.
	if s.Config.RunMode != "local" {
		urlPath = "http://localhost" + urlPath
	}
	httpReq, err := http.NewRequest(req.HTTPMethod, urlPath, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	return httpReq, nil
}
