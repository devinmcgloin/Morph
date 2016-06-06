package handler

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

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
				log.Fatal(err)
			}

			var buf bytes.Buffer
			var written int64
			written, err := buf.ReadFrom(infile)
			if err != nil {
				log.Fatal(err)
			}

			URL := dbase.UploadImageAWS(buf.Bytes(), written, hdr.Filename, "morph-content", "us-east-1")

			Img.URL = URL
			w.Write([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
		}
	}

	log.Printf("%v", Img)
	dbase.AddImg(Img, dbase.DB)

}

func demoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseMultipartForm(32 << 20)

	var (
		err error
	)

	for _, fheaders := range r.MultipartForm.File {
		for _, hdr := range fheaders {
			// open uploaded
			var infile multipart.File
			if infile, err = hdr.Open(); nil != err {
				log.Fatal(err)
			}
			// open destination
			var outfile *os.File
			if outfile, err = os.Create("./content/" + hdr.Filename); nil != err {
				log.Fatal(err)
			}
			// 32K buffer copy
			var written int64
			if written, err = io.Copy(outfile, infile); nil != err {
				log.Fatal(err)
			}
			w.Write([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
		}
	}
}
