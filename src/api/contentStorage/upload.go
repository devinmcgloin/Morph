package contentStorage

import (
	"bytes"
	"log"
	"math"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/devinmcgloin/sprioc/src/spriocError"
)

var mediaTypeOptions = []string{"image/jp2", "image/png", "image/tiff", "image/bmp"}

// ImageProcessor manages uploading the original file to aws.
// TODO: In the future it should also spin off worker threads to
// handle compression, and rendering other sizes for the image.
func ProccessImage(infile []byte, size int, id bson.ObjectId) error {

	var err error

	buf := bytes.NewBuffer(infile)

	contentType := http.DetectContentType(infile)
	in := in(contentType, mediaTypeOptions)
	if !in {
		log.Println(contentType)
		return spriocError.New(err, "Unsupported Media Type", 415)
	}

	// Check if image is smaller than 18 MB.
	if len(infile) > int(math.Exp2(20)*18) {
		return spriocError.New(nil, "Payload Too Large", 413)
	}

	err = UploadImageAWS(buf.Bytes(), int64(size), string(id), "sprioc", "us-east-1")
	if err != nil {
		return spriocError.New(err, "Error while uploading image", 500)
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
