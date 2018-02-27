package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"

	"github.com/fokal/fokal/pkg/metadata"
	"github.com/rwcarlsen/goexif/exif"
)

func main() {

	var in, out string

	flag.StringVar(&in, "in", "", "input file")
	flag.StringVar(&out, "out", "", "output file")

	flag.Parse()

	fmt.Printf("Reading from %s Writing to: %s\n", in, out)

	file, err := ioutil.ReadFile(in)
	if err != nil {
		log.Fatal("Invalid iamge path")
	}

	meta, err := metadata.GetExif(bytes.NewBuffer(file))
	if err != nil {
		log.Fatal("unable to load exif data")
	}

	var orientation int

	tag, err := meta.Get(exif.Orientation)
	if err == nil {
		orientation, err = tag.Int(0)
		if err != nil {
			log.Fatal("Invalid Orientation")
		}
	}

	uploadedImage, _, err := image.Decode(bytes.NewBuffer(file))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Orientation: %d\n", orientation)
	rotatedImage := metadata.NormalizeOrientatation(uploadedImage, uint16(orientation))

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, rotatedImage, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(out, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal("Unable to write to file")
	}

}
