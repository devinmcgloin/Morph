package domain

import (
	"context"
	"image"
)

type VisionService interface {
	AnnotateImage(ctx context.Context, img image.Image) (*ImageAnnotation, error)
}

type ImageAnnotation struct {
	Labels          []Label
	Safe            bool
	ColorProperties []Color
	Landmark        []Landmark
}
