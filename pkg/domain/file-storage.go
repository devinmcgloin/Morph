package domain

import "bytes"

type StorageService interface {
	UploadImage(img *bytes.Buffer, path string) error
	DeleteImage(path string) error
}
