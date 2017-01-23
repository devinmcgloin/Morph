package metadata

import (
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/sprioc/composer/pkg/model"
	gj "github.com/sprioc/geojson"
)

func GetExif(image io.Reader) (*exif.Exif, error) {
	exifDat, err := exif.Decode(image)
	if err != nil {
		return &exif.Exif{}, errors.New("Unable to parse exif")
	}

	return exifDat, nil
}

func GetMetadata(file io.Reader, img *model.Image) error {
	x, err := GetExif(file)
	if err != nil {
		return err
	}

	lat, lon, err := x.LatLong()
	if err != nil {
		img.Location = nil
	} else {
		point := gj.NewPoint(gj.Coordinate{gj.CoordType(lon), gj.CoordType(lat)})
		point.Type = "Point"
		img.Location = point
	}

	captureTime, err := x.DateTime()
	if err == nil {
		img.CaptureTime = captureTime.Unix()
	}

	//	Classic stats
	ExposureTime, err := x.Get(exif.ExposureTime)
	if err == nil {
		num, den, err := ExposureTime.Rat2(0)

		if err == nil {
			img.ExposureTime = strconv.FormatInt(num, 10) + "/" + strconv.FormatInt(den, 10)
		} else {
			log.Println(err)
		}
	}

	Aperture, err := x.Get(exif.ApertureValue)
	if err == nil {
		num, den, err := Aperture.Rat2(0)
		if err == nil {
			img.Aperture = strconv.FormatInt(num/den, 10)
		}
	}

	FocalLength, err := x.Get(exif.FocalLength)
	if err == nil {
		num, den, err := FocalLength.Rat2(0)
		if err == nil {
			img.FocalLength = strconv.FormatInt(num/den, 10)
		}
	}

	ISO, err := x.Get(exif.ISOSpeedRatings)
	if err == nil {
		short, err := ISO.Int(0)
		if err == nil {
			img.ISO = short
		}
	}

	// Make and model info
	Make, err := x.Get(exif.Make)
	if err == nil {
		str, err := Make.StringVal()
		if err == nil {
			img.Make = str
		}
	}

	Model, err := x.Get(exif.Model)
	if err == nil {
		str, err := Model.StringVal()
		if err == nil {
			img.Model = str
		}
	}

	LensMake, err := x.Get(exif.LensMake)
	if err == nil {
		str, err := LensMake.StringVal()
		if err == nil {
			img.LensMake = str
		}
	}

	LensModel, err := x.Get(exif.LensModel)
	if err == nil {
		str, err := LensModel.StringVal()
		if err == nil {
			img.LensModel = str
		}
	}

	// Setting fields in sources for orig image

	PixelXDimension, err := x.Get(exif.PixelXDimension)
	if err == nil {
		n, err := PixelXDimension.Int64(0)
		if err == nil {
			img.PixelXDimension = n
		}
	}

	PixelYDimension, err := x.Get(exif.PixelYDimension)
	if err == nil {
		n, err := PixelYDimension.Int64(0)
		if err == nil {
			img.PixelYDimension = n
		}
	}

	return nil
}
