package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

// mockS3Client is a mock implementation of the S3 client for testing.
type mockS3Client struct {
	GetObjectFunc func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func (m *mockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.GetObjectFunc(ctx, params, optFns...)
}

func TestS3Client_GetObject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockClient := &mockS3Client{
			GetObjectFunc: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte("test data"))),
				}, nil
			},
		}

		s3Svc := &s3Client{client: mockClient}
		body, err := s3Svc.GetObject(context.TODO(), "test-bucket", "test-key")
		assert.NoError(t, err)
		defer body.Close()

		data, err := io.ReadAll(body)
		assert.NoError(t, err)
		assert.Equal(t, "test data", string(data))
	})

	t.Run("error", func(t *testing.T) {
		mockClient := &mockS3Client{
			GetObjectFunc: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return nil, errors.New("s3 error")
			},
		}

		s3Svc := &s3Client{client: mockClient}
		_, err := s3Svc.GetObject(context.TODO(), "test-bucket", "test-key")
		assert.Error(t, err)
		assert.Equal(t, "s3 error", err.Error())
	})
}
