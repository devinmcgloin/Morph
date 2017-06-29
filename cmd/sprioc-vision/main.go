package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/vision/v1"
)

var visionService *vision.Service

func main() {
	var err error

	flags := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(flags)

	client := &http.Client{
		Transport: &transport.APIKey{Key: os.Getenv("GOOGLE_API_TOKEN")},
	}
	visionService, err = vision.New(client)
	if err != nil {
		log.Fatal(err)
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-image>\n", filepath.Base(os.Args[0]))
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(args[0]); err != nil {
		// Comes here if run() returns an error
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

}

func run(file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Construct a text request, encoding the image in base64.
	req := &vision.AnnotateImageRequest{
		// Apply image which is encoded by base64
		Image: &vision.Image{
			Content: base64.StdEncoding.EncodeToString(b),
		},
		// Apply features to indicate what type of image detection
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

	// Parse annotations from responses
	b, err = json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Print(err)
	}

	ioutil.WriteFile("./temp.json", b, 0644)
	return nil
}
