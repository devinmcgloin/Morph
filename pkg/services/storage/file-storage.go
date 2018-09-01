package storage

import (
	"context"
	"image"
)

type Service interface {
	UploadImage(ctx context.Context, img image.Image, path string) error
	DeleteImage(ctx context.Context, path string) error
}
