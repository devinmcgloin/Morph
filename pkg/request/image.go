package request

import (
	"net/http"

	"github.com/mholt/binding"
)

type PatchImageRequest struct {
	Tags []string `json:"tags" structs:"tags,omitempty"`

	Aperture     float64 `json:"aperture" structs:"aperture,omitempty"`
	ExposureTime string  `json:"exposure_time" structs:"exposure_time,omitempty"`
	FocalLength  int     `json:"focal_length" structs:"focal_length,omitempty"`
	ISO          int     `json:"iso" structs:"iso,omitempty"`

	Make      string `json:"make" structs:"make,omitempty"`
	Model     string `json:"model" structs:"model,omitempty"`
	LensMake  string `json:"lens_make" structs:"lens_make,omitempty"`
	LensModel string `json:"lens_model" structs:"lens_model,omitempty"`

	CaptureTime string `json:"capture_time" structs:"capture_time,omitempty"`

	Geo *GeoPatch `json:"geo" structs:"geo,omitempty"`
}

type GeoPatch struct {
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lng"`
	Description string  `json:"description"`
}

func (cf *PatchImageRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Tags:         "tags",
		&cf.Aperture:     "aperture",
		&cf.ExposureTime: "exposure_time",
		&cf.FocalLength:  "focal_length",
		&cf.ISO:          "iso",
		&cf.Make:         "make",
		&cf.Model:        "model",
		&cf.LensMake:     "lens_make",
		&cf.LensModel:    "lens_model",
		&cf.CaptureTime:  "capture_time",
	}
}
