package domain

import (
	"bytes"
	"context"
)

type StorageService interface {
	UploadImage(ctx context.Context, img *bytes.Buffer, mimeType, path string) error
	DeleteImage(ctx context.Context, path string) error
}
