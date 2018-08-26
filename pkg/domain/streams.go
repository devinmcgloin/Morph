package domain

import (
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
	StreamByID(id uint64) (*Stream, error)
	StreamsByCreator(userID uint64) (*[]Stream, error)
	CreateStream(creator uint64, title string) error
	SetDiscription(description string) error
	SetTitle(title string) error
}
