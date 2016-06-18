package model

import (
	"time"

	gj "github.com/kpawlik/geojson"
	"gopkg.in/mgo.v2/bson"
)

// Image contains all the proper information for rendering a single photo
type Image struct {
	ID           bson.ObjectId `bson:"_id"`
	ShortCode    string        `bson:"shortcode"`
	Title        string        `bson:"title,omitempty"`
	Desc         string        `bson:"desc,omitempty"`
	Aperture     int64         `bson:"aperture,omitempty"`
	ExposureTime int64         `bson:"exposure_time,omitempty"`
	FocalLength  int64         `bson:"focal_length,omitempty"`
	ISO          int64         `bson:"iso,omitempty"`
	Orientation  string        `bson:"orientation,omitempty"`
	CameraBody   string        `bson:"camera_body,omitempty"`
	Lens         string        `bson:"lens,omitempty"`
	Tags         []string      `bson:"tags,omitempty"`
	MachineTags  []string      `bson:"machine_tags,omitempty"`
	AlbumID      bson.ObjectId `bson:"album_id,omitempty"`
	CaptureTime  time.Time     `bson:"capture_time"`
	PublishTime  time.Time     `bson:"publish_time"`
	ImgDirection float64       `bson:"direction"`
	UserID       bson.ObjectId `bson:"user_id,omitempty"` // FIXME omitempty is temp for testing.
	Location     gj.Point      `bson:"location,omitempty"`
	Sources      []ImgSource   `bson:"sources"`
	Featured     bool          `bson:"featured,omitempty"`
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
	URL        string `bson:"url"`
	Resolution int64  `bson:"resolution,omitempty"`
	Width      int64  `bson:"width,omitempty"`
	Height     int64  `bson:"height,omitempty"`
	Size       string `bson:"size"`
	FileType   string `bson:"file_type"`
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
