package contentStorage

import (
	"bytes"
	"errors"
	"log"
	"math"
	"net/http"
	"strings"
)

var mediaTypeOptions = []string{"image/jp2", "image/jpeg", "image/png", "image/tiff", "image/bmp"}

// ProccessImage manages uploading the original file to aws.
func ProccessImage(infile []byte, size int, shortcode string, kind string) error {

	var err error

	buf := bytes.NewBuffer(infile)

	// TODO this does not match properly to the mediaTypeOptions
	contentType := http.DetectContentType(infile)
	in := in(contentType, mediaTypeOptions)
	if !in {
		log.Println(contentType)
		return errors.New("Unsupported Media Type")
	}

	// Check if image is smaller than 18 MB.
	if len(infile) > int(math.Exp2(20)*18) {
		return errors.New("Payload Too Large")
	}

	path := strings.Join([]string{"/images", kind, shortcode}, "/")

	err = UploadImageAWS(buf.Bytes(), int64(size), path, "sprioc", "us-east-1")
	if err != nil {
		return errors.New("Error while uploading image")
	}
	return nil
}

func in(contentType string, opts []string) bool {
	for _, opt := range opts {
		if strings.Compare(contentType, opt) == 0 {
			return true
		}
	}
	return false
}
