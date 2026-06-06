package storage

import (
	"context"
	"io"
)

type FileStorage interface {
	Upload(ctx context.Context, bucket, key string, reader io.Reader, size int64, contentType string) error
	Download(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket, key string) error
	Exists(ctx context.Context, bucket, key string) (bool, error)
}
