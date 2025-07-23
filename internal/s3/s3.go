package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Service defines the interface for interacting with S3.
type S3Service interface {
	GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)
}

// s3Client is an implementation of the S3Service interface.
type s3Client struct {
	client *s3.Client
}

// NewS3Service creates a new S3 service client.
func NewS3Service() (S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return &s3Client{client: client}, nil
}

// GetObject retrieves an object from S3.
func (s *s3Client) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}
