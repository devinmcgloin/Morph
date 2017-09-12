package request

import (
	"net/http"

	"github.com/mholt/binding"
)

type PatchImageRequest struct {
	Tags []string `structs:"tags,omitempty"`

	Aperture     float64 `structs:"aperture,omitempty"`
	ExposureTime string  `structs:"exposure_time,omitempty"`
	FocalLength  int     `structs:"focal_length,omitempty"`
	ISO          int     `structs:"iso,omitempty"`

	Make      string `structs:"make,omitempty"`
	Model     string `structs:"model,omitempty"`
	LensMake  string `json:"lens_make" structs:"lens_make,omitempty"`
	LensModel string `json:"lens_model" structs:"lens_model,omitempty"`

	CaptureTime string `json:"capture_time" structs:"capture_time,omitempty"`
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
