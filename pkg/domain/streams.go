package domain

import (
	"context"
	"time"
)

type Stream struct {
	ID          uint64
	Title       string
	Description string
	Creator     uint64

	UpdatedAt time.Time
	CreatedAt time.Time
}

//go:generate moq -out stream_service_runner.go . StreamService

type StreamService interface {
	StreamByID(ctx context.Context, id uint64) (*Stream, error)
	StreamsByCreator(ctx context.Context, userID uint64) (*[]Stream, error)
	CreateStream(ctx context.Context, creator uint64, title string) error
	SetDiscription(ctx context.Context, description string) error
	SetTitle(ctx context.Context, title string) error
}
