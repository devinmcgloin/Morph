package images

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/devinmcgloin/sprioc/src/api/AWS"
	"github.com/devinmcgloin/sprioc/src/api/auth"
	"github.com/devinmcgloin/sprioc/src/api/metadata"
	"github.com/devinmcgloin/sprioc/src/model"
	"github.com/devinmcgloin/sprioc/src/spriocError"

	"gopkg.in/mgo.v2/bson"
)

// UploadHandler manages uploading the original file to aws.
// TODO: In the future it should also spin off worker threads to
// handle compression, and rendering other sizes for the image.
func UploadHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseMultipartForm(32 << 20)

	var err error

	// TODO need to include shortext data here

	loggedIn, user := auth.CheckUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", 302)
		return nil
	}

	var imageSources []model.ImgSource
	imageSources = append(imageSources, model.ImgSource{Size: "orig"})

	image := model.Image{
		ID:          bson.NewObjectId(),
		UserID:      user.ID,
		PublishTime: time.Now(),
		CaptureTime: time.Now(),
		Sources:     imageSources,
		ShortCode:   mongo.GetNewImageShortCode(),
	}

	for _, fheaders := range r.MultipartForm.File {
		for _, hdr := range fheaders {

			// open uploaded
			var infile multipart.File
			infile, err = hdr.Open()
			if err != nil {
				return spriocError.New(err, "Error while reading in image", 500)
			}

			err := metadata.SetMetadata(infile, &image)
			if err != nil {
				return spriocError.New(err, "Error while reading image metadata", 500)
			}

			infile, err = hdr.Open()
			if err != nil {
				return spriocError.New(err, "Error while reading in image", 500)
			}

			var buf bytes.Buffer
			var written int64
			written, err = buf.ReadFrom(infile)
			if err != nil {
				return spriocError.New(err, "Error while reading in image", 500)
			}

			filename := fmt.Sprintf("%s_orig.jpg", image.ShortCode)
			image.Sources[0].FileType = http.DetectContentType(buf.Bytes())

			image.Sources[0].URL, err = AWS.UploadImageAWS(buf.Bytes(), written, filename, "morph-content", "us-east-1")
			if err != nil {
				return spriocError.New(err, "Error while uploading image", 500)
			}
		}
	}

	err = mongo.AddImg(image)
	if err != nil {
		return spriocError.New(err, "Error while adding image to db", 500)
	}

	err = mongo.AddUserImage(user.ID, image.ID)
	if err != nil {
		return spriocError.New(err, "Error while adding image to user", 500)

	}

	//TODO need to set image short title here
	newURL := fmt.Sprintf("/i/%s/", image.ShortCode)

	http.Redirect(w, r, newURL, 302)

	return nil
}
