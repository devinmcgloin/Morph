package endpoint

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/devinmcgloin/morph/src/api/AWS"
	"github.com/devinmcgloin/morph/src/api/SQL"
	"github.com/julienschmidt/httprouter"
)

// UploadHandler manages uploading the original file to aws.
// TODO: In the future it should also spin off worker threads to
// handle compression, and rendering other sizes for the image.
func UploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseMultipartForm(32 << 20)

	var err error

	img := SQL.Img{
		Title:       r.FormValue("Title"),
		Desc:        r.FormValue("Desc"),
		PublishTime: time.Now(),
		CaptureTime: time.Now(),
	}

	source := SQL.ImgSource{
		Size: "orig",
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

			source.URL, err = AWS.UploadImageAWS(buf.Bytes(), written, hdr.Filename, "morph-content", "us-east-1")
			if err != nil {
				log.Printf("Error while uploading image %s", err)
				http.Error(w, http.StatusText(500), 500)
				return
			}

		}
	}
	err = SQL.AddImg(img)
	if err != nil {
		log.Printf("Error while adding image to DB %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = SQL.AddSrc(source)
	if err != nil {
		log.Printf("Error while adding image to DB %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	http.Redirect(w, r, "/morph/upload", 302)

}
