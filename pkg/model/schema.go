package model

import (
	"googlemaps.github.io/maps"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/clarifai-go"
	gj "github.com/sprioc/geojson"
)

// TODO it would be good to have both public and private collections / images.

// Image contains all the proper information for rendering a single photo
type Image struct {
	ID          bson.ObjectId    `bson:"_id" json:"-"`
	MetaData    ImageMetaData    `bson:"metadata" json:"metadata"`
	Tags        []string         `bson:"tags" json:"tags"`
	MachineTags []string         `bson:"machine_tags" json:"machine_tags"`
	ColorTags   []clarifai.Color `bson:"color_tags" json:"color_tags"`
	PublishTime int64            `bson:"publish_time" json:"publish_time"`
	Sources     ImgSource        `bson:"sources_link" json:"source_link"`
}

type ImageMetaData struct {
	Aperture        string    `bson:"aperture" json:"aperture,omitempty"`
	ExposureTime    string    `bson:"exposure_time" json:"exposure_time,omitempty"`
	FocalLength     string    `bson:"focal_length" json:"focal_length,omitempty"`
	ISO             int       `bson:"iso" json:"iso,omitempty"`
	Make            string    `bson:"make" json:"make,omitempty"`
	Model           string    `bson:"model" json:"model,omitempty"`
	LensMake        string    `bson:"lens_make" json:"lens_make,omitempty"`
	LensModel       string    `bson:"lens_model" json:"lens_model,omitempty"`
	PixelXDimension int64     `bson:"pixel_xd" json:"pixel_xd,omitempty"`
	PixelYDimension int64     `bson:"pixel_yd" json:"pixel_yd,omitempty"`
	CaptureTime     int64     `bson:"capture_time" json:"capture_time,omitempty"`
	ImgDirection    float64   `bson:"direction" json:"direction,omitempty"`
	Location        *gj.Point `bson:"location" json:"location,omitempty"`
}

// ImgSource includes the information about the image itself.
type ImgSource struct {
	Thumb  string `bson:"thumb" json:"thumb"`
	Small  string `bson:"small" json:"small"`
	Medium string `bson:"medium" json:"medium"`
	Large  string `bson:"large" json:"large"`
	Raw    string `bson:"raw" json:"raw"`
}

type User struct {
	ID        bson.ObjectId `bson:"_id" json:"-"`
	Email     string        `bson:"email" json:"email"`
	Pass      string        `bson:"password" json:"-"`
	Salt      string        `bson:"salt" json:"-"`
	Name      string        `bson:"name" json:"name"`
	Bio       string        `bson:"bio" json:"bio,omitempty"`
	URL       string        `bson:"personal_site_link" json:"personal_site_link,omitempty"`
	Location  *gj.Feature   `bson:"location" json:"location,omitempty"`
	AvatarURL ImgSource     `bson:"avatar_link" json:"avatar_link"`
}

type Collection struct {
	ID       bson.ObjectId `bson:"_id" json:"-"`
	Desc     string        `bson:"desc" json:"desc,omitempty"`
	Title    string        `bson:"title" json:"title,omitempty"`
	Location *gj.Feature   `bson:"location" json:"location,omitempty"`
}

type Location struct {
	GoogleLoc maps.GeocodingResult `bson:"google_geo"`
	Bounds    gj.Polygon           `bson:"bounds"`
}
