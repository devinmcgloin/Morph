package metadata

import (
	"encoding/base64"
	"fmt"
	"io"

	"google.golang.org/api/vision/v1"
)

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
	fmt.Println(res)
	return nil

}
