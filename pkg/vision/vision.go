package vision

import (
	"encoding/base64"
	"log"

	"github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/clr/clr"
	"github.com/fokal/fokal/pkg/color"
	"github.com/fokal/fokal/pkg/model"
	"github.com/jmoiron/sqlx"

	"image"

	"bytes"
	"image/jpeg"

	"github.com/nfnt/resize"
	"google.golang.org/api/vision/v1"
)

type ImageResponse struct {
	Labels          []model.Label
	Safe            bool
	ColorProperties []model.Color
	Landmark        []model.Landmark
}

func AnnotateImage(errChan chan error, annotations chan ImageResponse, db *sqlx.DB, visionService *vision.Service, img image.Image) {

	m := resize.Resize(300, 0, img, resize.Bilinear)
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, m, nil)
	if err != nil {
		errChan <- err
		return
	}
	// Construct a text request, encoding the image in base64.

	req := &vision.AnnotateImageRequest{
		// Apply image which is encoded by base64
		Image: &vision.Image{
			Content: base64.StdEncoding.EncodeToString(buf.Bytes()),
		},
		// Apply features to indicate what type of image detection
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
		log.Println(err)
		errChan <- err
		return
	}

	r := res.Responses[0]
	rsp := ImageResponse{Safe: true}

	shade := color.RetrieveColorTable(db, color.Shade)
	specific := color.RetrieveColorTable(db, color.SpecificColor)

	for _, col := range r.ImagePropertiesAnnotation.DominantColors.Colors {
		sRGB := clr.RGB{
			R: uint8(col.Color.Red),
			G: uint8(col.Color.Green),
			B: uint8(col.Color.Blue)}

		h, s, v := sRGB.HSV()
		rsp.ColorProperties = append(rsp.ColorProperties, model.Color{
			SRGB:          sRGB,
			PixelFraction: col.PixelFraction,
			Score:         col.Score,
			Hex:           sRGB.Hex(),
			HSV: clr.HSV{
				H: h, S: s, V: v,
			},
			Shade:     sRGB.ColorName(shade),
			ColorName: sRGB.ColorName(specific),
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

	unique := make(map[string]bool, len(rsp.Labels))

	for _, label := range r.LabelAnnotations {
		if _, ok := unique[label.Description]; !ok {
			rsp.Labels = append(rsp.Labels, model.Label{
				Description: label.Description,
				Score:       label.Score,
			})
			unique[label.Description] = true
		}
	}

	for _, landmark := range r.LandmarkAnnotations {
		rsp.Landmark = append(rsp.Landmark, model.Landmark{
			Description: landmark.Description,
			Score:       landmark.Score,
			Location: postgis.PointS{
				SRID: 4326,
				X:    landmark.Locations[0].LatLng.Longitude,
				Y:    landmark.Locations[0].LatLng.Latitude,
			},
		})
	}
	errChan <- nil
	annotations <- rsp
}
