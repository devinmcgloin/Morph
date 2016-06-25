package metadata

import (
	"errors"
	"io"
	"log"
	"strconv"

	gj "github.com/kpawlik/geojson"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/sprioc/sprioc-core/pkg/model"
)

func GetExif(image io.Reader) (*exif.Exif, error) {
	exifDat, err := exif.Decode(image)
	if err != nil {
		return &exif.Exif{}, errors.New("Unable to parse exif")
	}

	return exifDat, nil
}

func GetMetadata(file io.Reader) (model.ImageMetaData, error) {
	x, err := GetExif(file)
	if err != nil {
		return model.ImageMetaData{}, nil
	}

	metaData := model.ImageMetaData{}

	// TODO: Could add altitude data here.
	lat, lon, err := x.LatLong()
	if err == nil {
		point := gj.NewPoint(gj.Coordinate{gj.CoordType(lon), gj.CoordType(lat)})
		point.Type = "Point"
		metaData.Location = *gj.NewFeature(point, nil, nil)
	}

	tmp, err := x.DateTime()
	if err == nil {
		metaData.CaptureTime = tmp
	}

	// TODO need to consider other represenations

	//	Classic stats
	ExposureTime, err := x.Get(exif.ExposureTime)
	if err == nil {
		num, den, err := ExposureTime.Rat2(0)

		if err == nil {
			metaData.ExposureTime = model.NewRatio(num, den, strconv.FormatInt(num, 10)+"/"+strconv.FormatInt(den, 10))
		} else {
			log.Println(err)
		}
	}

	Aperture, err := x.Get(exif.ApertureValue)
	if err == nil {
		num, den, err := Aperture.Rat2(0)
		if err == nil {
			metaData.Aperture = model.NewRatio(num, den, strconv.FormatInt(num/den, 10))
		}
	}

	FocalLength, err := x.Get(exif.FocalLength)
	if err == nil {
		num, den, err := FocalLength.Rat2(0)
		if err == nil {
			metaData.FocalLength = model.NewRatio(num, den, strconv.FormatInt(num/den, 10))
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

	return metaData, nil
}
