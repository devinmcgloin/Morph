package vision

import (
	"context"

	img "image"

	"github.com/fokal/fokal-core/pkg/services/image"
)

type Service interface {
	AnnotateImage(ctx context.Context, img img.Image) (*ImageAnnotation, error)
}

type ImageAnnotation struct {
	Labels          []image.Label
	Safe            bool
	ColorProperties []image.Color
	Landmark        []image.Landmark
}
