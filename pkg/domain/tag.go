package domain

import "context"

type Tag struct {
	ID          uint64 `json:"id"`
	Description string
}

//go:generate moq -out tag_service_runner.go . TagService

type TagService interface {
	TagByID(ctx context.Context, id uint64) (*Tag, error)
	TagByDescription(ctx context.Context, desc string) (*Tag, error)
	CreateTag(ctx context.Context, desc string) (*Tag, error)
	ImagesForTag(ctx context.Context, id uint64) (*[]Image, error)
	TagImage(ctx context.Context, id uint64, imageID uint64) error
	UnTagImage(ctx context.Context, id uint64, imageID uint64) error
}
