package api

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/devinmcgloin/morph/src/dbase"
	"github.com/devinmcgloin/morph/src/schema"
	"github.com/julienschmidt/httprouter"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseMultipartForm(32 << 20)

	var err error

	Img := schema.Img{
		Title:    r.FormValue("Title"),
		Desc:     r.FormValue("Desc"),
		Category: r.FormValue("Category"),
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

			var URL string
			URL, err = dbase.UploadImageAWS(buf.Bytes(), written, hdr.Filename, "morph-content", "us-east-1")
			if err != nil {
				log.Printf("Error while uploading image %s", err)
				http.Error(w, http.StatusText(500), 500)
				return
			}
			Img.URL = URL
		}
	}

	log.Printf("%v", Img)
	err = dbase.AddImg(Img)
	if err != nil {
		log.Printf("Error while adding image to DB %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	http.Redirect(w, r, "/morph/manage", 302)

}
