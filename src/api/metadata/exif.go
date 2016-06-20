package metadata

import (
	"io"
	"log"
	"strconv"

	"github.com/devinmcgloin/morph/src/model"
	"github.com/devinmcgloin/morph/src/morphError"
	gj "github.com/kpawlik/geojson"
	"github.com/rwcarlsen/goexif/exif"
)

func GetExif(image io.Reader) (*exif.Exif, error) {
	exifDat, err := exif.Decode(image)
	if err != nil {
		return &exif.Exif{}, morphError.New(err, "Unable to parse exif", 523)
	}

	return exifDat, nil
}

func SetMetadata(file io.Reader, image *model.Image) error {
	x, err := GetExif(file)
	if err != nil {
		return nil
	}

	// TODO: Could add altitude data here.
	lat, lon, err := x.LatLong()
	if err == nil {
		point := gj.NewPoint(gj.Coordinate{gj.CoordType(lon), gj.CoordType(lat)})
		point.Type = "Point"
		image.Location = *point
	}

	tmp, err := x.DateTime()
	if err == nil {
		image.CaptureTime = tmp
	}

	// TODO need to consider other represenations

	//	Classic stats
	ExposureTime, err := x.Get(exif.ExposureTime)
	if err == nil {
		num, den, err := ExposureTime.Rat2(0)

		if err == nil {
			image.ExposureTime = model.NewRatio(num, den, strconv.FormatInt(num, 10)+"/"+strconv.FormatInt(den, 10))
		} else {
			log.Println(err)
		}
	}

	Aperture, err := x.Get(exif.ApertureValue)
	if err == nil {
		num, den, err := Aperture.Rat2(0)
		if err == nil {
			image.Aperture = model.NewRatio(num, den, strconv.FormatInt(num/den, 10))
		}
	}

	FocalLength, err := x.Get(exif.FocalLength)
	if err == nil {
		num, den, err := FocalLength.Rat2(0)
		if err == nil {
			image.FocalLength = model.NewRatio(num, den, strconv.FormatInt(num/den, 10))
		}
	}

	ISO, err := x.Get(exif.ISOSpeedRatings)
	if err == nil {
		short, err := ISO.Int(0)
		if err == nil {
			image.ISO = short
		}
	}

	// Make and model info
	Make, err := x.Get(exif.Make)
	if err == nil {
		str, err := Make.StringVal()
		if err == nil {
			image.Make = str
		}
	}

	Model, err := x.Get(exif.Model)
	if err == nil {
		str, err := Model.StringVal()
		if err == nil {
			image.Model = str
		}
	}

	LensMake, err := x.Get(exif.LensMake)
	if err == nil {
		str, err := LensMake.StringVal()
		if err == nil {
			image.LensMake = str
		}
	}

	LensModel, err := x.Get(exif.LensModel)
	if err == nil {
		str, err := LensModel.StringVal()
		if err == nil {
			image.LensModel = str
		}
	}

	// Setting fields in sources for orig image

	PixelXDimension, err := x.Get(exif.PixelXDimension)
	if err == nil {
		n, err := PixelXDimension.Int64(0)
		if err == nil {
			image.Sources[0].PixelXDimension = n
		}
	}

	PixelYDimension, err := x.Get(exif.PixelYDimension)
	if err == nil {
		n, err := PixelYDimension.Int64(0)
		if err == nil {
			image.Sources[0].PixelYDimension = n
		}
	}

	return nil
}
