package model

import (
	"fmt"
	"time"

	gj "github.com/kpawlik/geojson"
	"gopkg.in/mgo.v2/bson"
)

// Image contains all the proper information for rendering a single photo
type Image struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	ShortCode    string        `bson:"shortcode" json:"shortcode"`
	Aperture     Ratio         `bson:"aperture,omitempty" json:"aperture,omitempty"`
	ExposureTime Ratio         `bson:"exposure_time,omitempty" json:"exposure_time,omitempty"`
	FocalLength  Ratio         `bson:"focal_length,omitempty" json:"focal_length,omitempty"`
	ISO          int           `bson:"iso,omitempty" json:"iso,omitempty"`
	Orientation  string        `bson:"orientation,omitempty" json:"orientation,omitempty"`
	Make         string        `bson:"make,omitempty" json:"make,omitempty"`
	Model        string        `bson:"model,omitempty" json:"model,omitempty"`
	LensMake     string        `bson:"lens_make,omitempty" json:"lens_make,omitempty"`
	LensModel    string        `bson:"lens_model,omitempty" json:"lens_model,omitempty"`
	Tags         []string      `bson:"tags,omitempty" json:"tags,omitempty"`
	MachineTags  []string      `bson:"machine_tags,omitempty" json:"machine_tags,omitempty"`
	AlbumID      bson.ObjectId `bson:"album_id,omitempty" json:"album_id,omitempty"`
	CaptureTime  time.Time     `bson:"capture_time" json:"capture_time"`
	PublishTime  time.Time     `bson:"publish_time" json:"publish_time"`
	ImgDirection float64       `bson:"direction,omitempty" json:"direction,omitempty"`
	UserID       bson.ObjectId `bson:"user_id" json:"user_id"`
	Location     gj.Point      `bson:"location,omitempty" json:"location,omitempty"`
	Sources      []ImgSource   `bson:"sources" json:"sources"`
	Featured     bool          `bson:"featured,omitempty" json:"featured,omitempty"`
}

// Location fields that are currently unused as gj.Point seems a better fit.
type Location struct {
	Latitude  float64 `bson:"lat"`
	Longitude float64 `bson:"lon"`
	Desc      string  `bson:"desc"`
	ShortText string  `bson:"short_text"`
}

// ImgSource includes the information about the image itself.
// Size indicates how large the image is.
type ImgSource struct {
	URL             string `bson:"url" json:"url"`
	PixelXDimension int64  `bson:"pixel_xd,omitempty" json:"pixel_xd,omitempty"`
	PixelYDimension int64  `bson:"pixel_yd,omitempty" json:"pixel_yd,omitempty"`
	Size            string `bson:"size" json:"size"`
	FileType        string `bson:"file_type,omitempty" json:"file_type,omitempty"`
}

type Ratio struct {
	Num int64  `bson:"-" json:"-"`
	Den int64  `bson:"-" json:"-"`
	Rep string `bson:"rep" json:"rep"`
}

func (r Ratio) Format() string {
	return fmt.Sprintf("%d", (r.Num / r.Den))
}

// NewRatio is a shortcut to make ratio types
func NewRatio(num, den int64, rep string) Ratio {
	return Ratio{Num: num, Den: den, Rep: rep}
}

// User collection schema
type User struct {
	ID         bson.ObjectId   `bson:"_id"`
	Images     []bson.ObjectId `bson:"images"`
	UserName   string          `bson:"username"`
	Email      string          `bson:"email"`
	Name       string          `bson:"name"`
	Bio        string          `bson:"bio,omitempty"`
	Location   gj.Point        `bson:"loc"`
	AvatarURL  string          `bson:"avatar_url"`
	Provider   string          `bson:"provider"`
	ProviderID string          `bson:"provider_id"`
}

// Album collection schema
type Album struct {
	ID        bson.ObjectId   `bson:"_id"`
	UserID    bson.ObjectId   `bson:"user_id"`
	Images    []bson.ObjectId `bson:"images"`
	Desc      string          `bson:"desc"`
	Title     string          `bson:"title"`
	ViewType  string          `bson:"view_type"`
	UserName  string          `bson:"username"`
	ShortCode string          `bson:"shortcode"`
}
