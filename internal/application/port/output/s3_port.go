package output

import (
	"context"
	"io"
)

// S3Port defines the interface for S3 storage operations.
type S3Port interface {
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	PutObject(ctx context.Context, key string, body io.Reader, contentType string) error
	DeleteObject(ctx context.Context, key string) error
}
