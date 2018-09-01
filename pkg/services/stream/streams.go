package stream

import (
	"context"
	"time"

	"github.com/fokal/fokal-core/pkg/services/image"
)

type Stream struct {
	ID          uint64
	Title       string
	Description *string
	Creator     uint64

	UpdatedAt time.Time
	CreatedAt time.Time
}

//go:generate moq -out stream_service_runner.go . StreamService

type Service interface {
	StreamByID(ctx context.Context, id uint64) (*Stream, error)
	StreamsByCreator(ctx context.Context, userID uint64) (*[]Stream, error)
	CreateStream(ctx context.Context, creator uint64, title string) error
	SetDescription(ctx context.Context, id uint64, description string) error
	SetTitle(ctx context.Context, id uint64, title string) error
	AddImage(ctx context.Context, id, imageID uint64) error
	RemoveImage(ctx context.Context, id, imageID uint64) error
	Images(ctx context.Context, id uint64) (*[]image.Image, error)
}
