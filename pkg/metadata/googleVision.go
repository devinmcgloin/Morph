package metadata

import (
	"encoding/base64"
	"io"

	"github.com/sprioc/clr/clr"
	"github.com/sprioc/composer/pkg/model"
	gj "github.com/sprioc/geojson"
	"google.golang.org/api/vision/v1"
)

type ImageResponse struct {
	Labels          []model.Label
	Safe            bool
	ColorProperties []model.Color
	Landmark        []model.Landmark
}

func AnnotateImage(file io.Reader) (ImageResponse, error) {
	var b []byte
	_, err := file.Read(b)
	if err != nil {
		return ImageResponse{}, err
	}

	req := &vision.AnnotateImageRequest{
		Image: &vision.Image{
			Content: base64.StdEncoding.EncodeToString(b),
		},

		Features: []*vision.Feature{
			{Type: "SAFE_SEARCH_DETECTION"},
			{Type: "LANDMARK_DETECTION"},
			{Type: "IMAGE_PROPERTIES"},
			{Type: "LABEL_DETECTION"},
		},
	}

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}

	res, err := visionService.Images.Annotate(batch).Do()
	if err != nil {
		return ImageResponse{}, err
	}

	r := res.Responses[0]
	rsp := ImageResponse{Safe: true}
	for _, col := range r.ImagePropertiesAnnotation.DominantColors.Colors {
		sRGB := clr.RGB{
			R: uint8(255 * col.Color.Red),
			G: uint8(255 * col.Color.Green),
			B: uint8(255 * col.Color.Blue)}
		h, s, v := sRGB.HSV()
		rsp.ColorProperties = append(rsp.ColorProperties, model.Color{
			SRGB:          sRGB,
			PixelFraction: col.PixelFraction,
			Score:         col.Score,
			Hex:           sRGB.Hex(),
			HSV: clr.HSV{
				H: h, S: s, V: v,
			},
			Shade:     sRGB.Shade(),
			ColorName: sRGB.ColorName(),
		})
	}

	for _, likelihood := range []string{"POSSIBLE", "LIKELY", "VERY_LIKELY"} {
		if r.SafeSearchAnnotation.Adult == likelihood {
			rsp.Safe = false
		}
		if r.SafeSearchAnnotation.Violence == likelihood {
			rsp.Safe = false
		}
		if r.SafeSearchAnnotation.Medical == likelihood {
			rsp.Safe = false
		}
		if r.SafeSearchAnnotation.Spoof == likelihood {
			rsp.Safe = false
		}
	}

	for _, label := range r.LabelAnnotations {
		rsp.Labels = append(rsp.Labels, model.Label{
			Description: label.Description,
			Score:       label.Score,
		})
	}

	for _, landmark := range r.LandmarkAnnotations {
		rsp.Landmark = append(rsp.Landmark, model.Landmark{
			Description: landmark.Description,
			Score:       landmark.Score,
			Location: gj.Point{
				Type: "point",
				Coordinates: gj.Coordinate{
					gj.Coord(landmark.Locations[0].LatLng.Latitude),
					gj.Coord(landmark.Locations[0].LatLng.Longitude),
				},
			},
		})
	}

	return rsp, nil
}
