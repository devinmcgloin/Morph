package image

import (
	"log"
	"os"

	gj "github.com/kpawlik/geojson"
	"github.com/rwcarlsen/goexif/exif"
)

func GetExif(image *os.File) (*exif.Exif, *gj.Point, error) {
	exifDat, err := exif.Decode(image)
	if err != nil {
		log.Println(err)
		return &exif.Exif{}, &gj.Point{}, err
	}

	lat, lon, err := exifDat.LatLong()
	if err != nil {
		log.Println(err)
		return &exif.Exif{}, &gj.Point{}, err
	}
	p := gj.NewPoint(gj.Coordinate{gj.CoordType(lon), gj.CoordType(lat)})
	p.Type = "Point"
	return exifDat, p, nil
}
