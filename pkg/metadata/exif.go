package metadata

import (
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/sprioc/conductor/pkg/model"
	gj "github.com/sprioc/geojson"
)

func GetExif(image io.Reader) (*exif.Exif, error) {
	exifDat, err := exif.Decode(image)
	if err != nil {
		return &exif.Exif{}, errors.New("Unable to parse exif")
	}

	return exifDat, nil
}

func GetMetadata(file io.Reader, image *model.Image) {
	x, err := GetExif(file)
	if err != nil {
		return
	}

	metaData := model.ImageMetaData{}

	lat, lon, err := x.LatLong()
	if err == nil {
		point := gj.NewPoint(gj.Coordinate{gj.CoordType(lon), gj.CoordType(lat)})
		point.Type = "Point"
		metaData.Location = *point
	}

	tmp, err := x.DateTime()
	if err == nil {
		metaData.CaptureTime = tmp
	}

	//	Classic stats
	ExposureTime, err := x.Get(exif.ExposureTime)
	if err == nil {
		num, den, err := ExposureTime.Rat2(0)

		if err == nil {
			metaData.ExposureTime = strconv.FormatInt(num, 10) + "/" + strconv.FormatInt(den, 10)
		} else {
			log.Println(err)
		}
	}

	Aperture, err := x.Get(exif.ApertureValue)
	if err == nil {
		num, den, err := Aperture.Rat2(0)
		if err == nil {
			metaData.Aperture = strconv.FormatInt(num/den, 10)
		}
	}

	FocalLength, err := x.Get(exif.FocalLength)
	if err == nil {
		num, den, err := FocalLength.Rat2(0)
		if err == nil {
			metaData.FocalLength = strconv.FormatInt(num/den, 10)
		}
	}

	ISO, err := x.Get(exif.ISOSpeedRatings)
	if err == nil {
		short, err := ISO.Int(0)
		if err == nil {
			metaData.ISO = short
		}
	}

	// Make and model info
	Make, err := x.Get(exif.Make)
	if err == nil {
		str, err := Make.StringVal()
		if err == nil {
			metaData.Make = str
		}
	}

	Model, err := x.Get(exif.Model)
	if err == nil {
		str, err := Model.StringVal()
		if err == nil {
			metaData.Model = str
		}
	}

	LensMake, err := x.Get(exif.LensMake)
	if err == nil {
		str, err := LensMake.StringVal()
		if err == nil {
			metaData.LensMake = str
		}
	}

	LensModel, err := x.Get(exif.LensModel)
	if err == nil {
		str, err := LensModel.StringVal()
		if err == nil {
			metaData.LensModel = str
		}
	}

	// Setting fields in sources for orig image

	PixelXDimension, err := x.Get(exif.PixelXDimension)
	if err == nil {
		n, err := PixelXDimension.Int64(0)
		if err == nil {
			metaData.PixelXDimension = n
		}
	}

	PixelYDimension, err := x.Get(exif.PixelYDimension)
	if err == nil {
		n, err := PixelYDimension.Int64(0)
		if err == nil {
			metaData.PixelYDimension = n
		}
	}

	log.Println(metaData.Location)

	image.MetaData = metaData
}
