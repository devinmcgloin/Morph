package upload

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"log"
	"strings"
)

var mediaTypeOptions = []string{"jp2", "jpeg", "png", "tiff", "bmp"}

// ProccessImage manages uploading the original file to aws.
func ProccessImage(errChan chan error, img image.Image, format string, shortcode string, kind string) {

	var err error

	// TODO this does not match properly to the mediaTypeOptions
	in := in(format, mediaTypeOptions)
	if !in {
		log.Println(format)
		errChan <- errors.New("Unsupported Media Type")
		return
	}

	path := strings.Join([]string{kind, shortcode}, "/")
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		errChan <- err
		return
	}
	err = UploadImageAWS(buf, format, path, "images-fokal", "us-west-1")
	if err != nil {
		errChan <- errors.New("Error while uploading image")
		return
	}
	errChan <- nil
	return
}

func in(contentType string, opts []string) bool {
	for _, opt := range opts {
		if strings.Compare(contentType, opt) == 0 {
			return true
		}
	}
	return false
}
