package domain

import (
	"context"
	"image"
)

type StorageService interface {
	UploadImage(ctx context.Context, img image.Image, path string) error
	DeleteImage(ctx context.Context, path string) error
}
