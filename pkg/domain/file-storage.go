package domain

import "bytes"

type StorageService interface {
	UploadImage(img *bytes.Buffer, mimeType, path string) error
	DeleteImage(path string) error
}
