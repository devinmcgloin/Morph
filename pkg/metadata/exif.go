package metadata

import (
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/rwcarlsen/goexif/exif"
)

func GetExif(image io.Reader) (*exif.Exif, error) {
	exifDat, err := exif.Decode(image)
	if err != nil {
		return &exif.Exif{}, errors.New("Unable to parse exif")
	}

	return exifDat, nil
}

func GetMetadata(file io.Reader) (metadata model.ImageMetadata, err error) {
	x, err := GetExif(file)
	if err != nil {
		return
	}

	lat, lon, err := x.LatLong()
	if err != nil {
		metadata.Location = nil
	} else {
		point := &postgis.PointS{SRID: 4326,
			X: lon,
			Y: lat}

		metadata.Location = point
	}

	dirTag, err := x.Get(exif.GPSImgDirection)
	if err == nil {
		dir, err := dirTag.Float(0)
		if err == nil {
			metadata.ImageDirection = &dir
		}
	}

	captureTime, err := x.DateTime()
	if err == nil {
		metadata.CaptureTime = &captureTime
	}

	//	Classic stats
	ExposureTime, err := x.Get(exif.ExposureTime)
	if err == nil {
		num, den, err := ExposureTime.Rat2(0)

		if err == nil {
			et := strconv.FormatInt(num, 10) + "/" + strconv.FormatInt(den, 10)
			metadata.ExposureTime = &et
		} else {
			log.Println(err)
		}
	}

	Aperture, err := x.Get(exif.ApertureValue)
	if err == nil {
		num, den, err := Aperture.Rat2(0)
		if err == nil {
			a := strconv.FormatFloat(float64(num)/float64(den), 'f', 1, 64)
			metadata.Aperture = &a
		}
	}

	FocalLength, err := x.Get(exif.FocalLength)
	if err == nil {
		num, den, err := FocalLength.Rat2(0)
		if err == nil {
			fl := strconv.FormatInt(num/den, 10)
			metadata.FocalLength = &fl
		}
	}

	ISO, err := x.Get(exif.ISOSpeedRatings)
	if err == nil {
		short, err := ISO.Int(0)
		if err == nil {
			metadata.ISO = &short
		}
	}

	// Make and model info
	Make, err := x.Get(exif.Make)
	if err == nil {
		str, err := Make.StringVal()
		str = strings.TrimSpace(str)
		if err == nil {
			metadata.Make = &str
		}
	}

	Model, err := x.Get(exif.Model)
	if err == nil {
		str, err := Model.StringVal()
		str = strings.TrimSpace(str)
		if err == nil {
			metadata.Model = &str
		}
	}

	LensMake, err := x.Get(exif.LensMake)
	if err == nil {
		str, err := LensMake.StringVal()
		str = strings.TrimSpace(str)
		if err == nil {
			metadata.LensMake = &str
		}
	}

	LensModel, err := x.Get(exif.LensModel)
	if err == nil {
		str, err := LensModel.StringVal()
		str = strings.TrimSpace(str)
		if err == nil {
			metadata.LensModel = &str
		}
	}

	// Setting fields in sources for orig image

	PixelXDimension, err := x.Get(exif.PixelXDimension)
	if err == nil {
		n, err := PixelXDimension.Int64(0)
		if err == nil {
			metadata.PixelXDimension = &n
		}
	}

	PixelYDimension, err := x.Get(exif.PixelYDimension)
	if err == nil {
		n, err := PixelYDimension.Int64(0)
		if err == nil {
			metadata.PixelYDimension = &n
		}
	}

	return metadata, nil
}
