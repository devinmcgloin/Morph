package metadata

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/vision/v1"
)

var visionService *vision.Service

func init() {
	client := &http.Client{
		Transport: &transport.APIKey{Key: os.Getenv("GOOGLE_API_KEY")},
	}
	visionService, err := vision.New(client)
	if err != nil {
		log.Fatal(err)
	}
}

func GetResponse(file io.Reader) error {
	var b []byte
	_, err := file.Read(b)
	if err != nil {
		return err
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
		return err
	}
	return nil

}
