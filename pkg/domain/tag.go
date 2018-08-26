package domain

type Tag struct {
	ID          string `json:"id"`
	Description string
}

//go:generate moq -out tag_service_runner.go . TagService

type TagService interface {
	TagByID(id string) (*Tag, error)
	CreateTag(id string) error
	ImagesForTag(id string) (*[]Image, error)
	TagImage(id string, imageID uint64) error
	UnTagImage(id string, imageID uint64) error
}
