package endpoint

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/devinmcgloin/morph/src/api/AWS"
	"github.com/devinmcgloin/morph/src/model"
	"github.com/julienschmidt/httprouter"
)

// UploadHandler manages uploading the original file to aws.
// TODO: In the future it should also spin off worker threads to
// handle compression, and rendering other sizes for the image.
func UploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseMultipartForm(32 << 20)

	var err error

	// TODO need to include shortext data here

	var imageSources []model.ImgSource
	imageSources = append(imageSources, model.ImgSource{Size: "orig"})

	ID := bson.NewObjectId()
	image := model.Image{
		ID:          ID,
		PublishTime: time.Now(),
		CaptureTime: time.Now(),
		Sources:     imageSources,
		ShortTitle:  mongo.GetShortTitle(),
	}

	for _, fheaders := range r.MultipartForm.File {
		for _, hdr := range fheaders {

			// open uploaded
			var infile multipart.File
			infile, err = hdr.Open()

			if err != nil {
				log.Printf("Error while reading in image %s", err)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			var buf bytes.Buffer
			var written int64
			written, err = buf.ReadFrom(infile)
			if err != nil {
				log.Printf("Error while reading in image %s", err)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			filename := fmt.Sprintf("%s_orig.jpg", image.ShortTitle)
			log.Printf("Filename = %s", filename)

			image.Sources[0].URL, err = AWS.UploadImageAWS(buf.Bytes(), written, filename, "morph-content", "us-east-1")
			if err != nil {
				log.Printf("Error while uploading image %s", err)
				http.Error(w, http.StatusText(500), 500)
				return
			}

		}
	}

	log.Println(image.Sources[0].URL)

	err = mongo.AddImg(image)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//TODO need to set image short title here
	newURL := fmt.Sprintf("/i/%s/edit", image.ShortTitle)

	log.Println(newURL)

	http.Redirect(w, r, newURL, 302)

}
