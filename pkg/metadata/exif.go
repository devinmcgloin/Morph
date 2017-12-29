package metadata

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"io"

	"github.com/cridenour/go-postgis"
	"github.com/fokal/fokal/pkg/model"
	"github.com/rwcarlsen/goexif/exif"
)

func GetExif(img io.Reader) (*exif.Exif, error) {
	exifDat, err := exif.Decode(img)
	if err != nil {
		return &exif.Exif{}, errors.New("Unable to parse exif")
	}

	return exifDat, nil
}

func GetMetadata(errChan chan error, metaChan chan model.ImageMetadata, img io.Reader) {
	meta := model.ImageMetadata{Location: &model.Location{}}
	x, err := GetExif(img)
	if err != nil {
		errChan <- err
		return
	}

	lat, lng, err := x.LatLong()
	if err != nil {
		meta.Location = nil
	} else {
		point := &postgis.PointS{SRID: 4326,
			X: lng,
			Y: lat}

		meta.Location.Point = point
	}

	dirTag, err := x.Get(exif.GPSImgDirection)
	if err == nil {
		dir, err := dirTag.Float(0)
		if err == nil {
			meta.Location.ImageDirection = &dir
		}
	}

	captureTime, err := x.DateTime()
	if err == nil {
		meta.CaptureTime = &captureTime
	}

	//	Classic stats
	ExposureTime, err := x.Get(exif.ExposureTime)
	if err == nil {
		num, den, err := ExposureTime.Rat2(0)

		if err == nil {
			et := strconv.FormatInt(num, 10) + "/" + strconv.FormatInt(den, 10)
			meta.ExposureTime = &et
		} else {
			log.Println(err)
		}
	}

	Aperture, err := x.Get(exif.FNumber)
	if err == nil {
		num, den, err := Aperture.Rat2(0)
		if err == nil {
			a := Round(float64(num)/float64(den), 0.1)
			meta.Aperture = &a
		}
	}

	FocalLength, err := x.Get(exif.FocalLength)
	if err == nil {
		num, den, err := FocalLength.Rat2(0)
		if err == nil {
			f := int(num / den)
			meta.FocalLength = &f
		}
	}

	ISO, err := x.Get(exif.ISOSpeedRatings)
	if err == nil {
		short, err := ISO.Int(0)
		if err == nil {
			meta.ISO = &short
		}
	}

	// Make and model info
	Make, err := x.Get(exif.Make)
	if err == nil {
		str, err := Make.StringVal()
		str = strings.TrimSpace(str)
		if err == nil {
			meta.Make = &str
		}
	}

	Model, err := x.Get(exif.Model)
	if err == nil {
		str, err := Model.StringVal()
		str = strings.TrimSpace(str)
		if err == nil {
			meta.Model = &str
		}
	}

	LensMake, err := x.Get(exif.LensMake)
	if err == nil {
		str, err := LensMake.StringVal()
		str = strings.TrimSpace(str)
		if err == nil {
			meta.LensMake = &str
		}
	}

	LensModel, err := x.Get(exif.LensModel)
	if err == nil {
		str, err := LensModel.StringVal()
		str = strings.TrimSpace(str)
		if err == nil {
			meta.LensModel = &str
		}
	}

	// Setting fields in sources for orig image

	PixelXDimension, err := x.Get(exif.PixelXDimension)
	if err == nil {
		n, err := PixelXDimension.Int64(0)
		if err == nil {
			meta.PixelXDimension = &n
		}
	}

	PixelYDimension, err := x.Get(exif.PixelYDimension)
	if err == nil {
		n, err := PixelYDimension.Int64(0)
		if err == nil {
			meta.PixelYDimension = &n
		}
	}

	metaChan <- meta
	errChan <- nil
}

func Round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}
