package client

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
)

// S3Client implements the S3Port interface using AWS SDK v2.
type S3Client struct {
	client *s3.Client
	bucket string
}

// NewS3Client creates a new S3Client.
func NewS3Client(region, accessKey, secretKey, bucket string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Client{
		client: client,
		bucket: bucket,
	}, nil
}

// GetObject retrieves an object from S3.
func (c *S3Client) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}
	return result.Body, nil
}

// PutObject uploads an object to S3.
func (c *S3Client) PutObject(ctx context.Context, key string, body io.Reader, contentType string) error {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to put object to S3: %w", err)
	}
	return nil
}

// DeleteObject removes an object from S3.
func (c *S3Client) DeleteObject(ctx context.Context, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}
	return nil
}

var _ output.S3Port = (*S3Client)(nil)
