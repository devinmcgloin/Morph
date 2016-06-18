package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Image contains all the proper information for rendering a single photo
type Image struct {
	ID           bson.ObjectId `bson:"_id"`
	Slug         string        `bson:"slug"`
	Title        string        `bson:"title"`
	Desc         string        `bson:"desc"`
	Aperture     int64         `bson:"aperture"`
	ExposureTime int64         `bson:"exposure_time"`
	FocalLength  int64         `bson:"focal_length"`
	ISO          int64         `bson:"iso"`
	Orientation  string        `bson:"orientation"`
	CameraBody   string        `bson:"camera_body"`
	Lens         string        `bson:"lens"`
	Tags         []string      `bson:"tags"`
	AlbumID      bson.ObjectId `bson:"album_id"`
	CaptureTime  time.Time     `bson:"capture_time"`
	PublishTime  time.Time     `bson:"publish_time"`
	ImgDirection float64       `bson:"direction"`
	UserID       bson.ObjectId `bson:"user_id"`
	Location     Location      `bson:"location"`
	Sources      []ImgSource   `bson:"sources"`
}

// Location fields that get embedded into image and user
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
	Resolution int64  `bson:"resolution"`
	Width      int64  `bson:"width"`
	Height     int64  `bson:"height"`
	Size       string `bson:"size"`
	FileType   string `bson:"file_type"`
}

// User collection schema
type User struct {
	ID        bson.ObjectId   `bson:"_id"`
	Images    []bson.ObjectId `bson:"images"`
	UserName  string          `bson:"username"`
	Email     string          `bson:"email"`
	FirstName string          `bson:"first_name"`
	LastName  string          `bson:"last_name"`
	Bio       string          `bson:"bio"`
	Location  Location        `bson:"loc"`
	AvatarURL string          `bson:"avatar_url"`
}

// Album collection schema
type Album struct {
	ID       bson.ObjectId   `bson:"_id"`
	UserID   bson.ObjectId   `bson:"user_id"`
	Images   []bson.ObjectId `bson:"images"`
	Desc     string          `bson:"desc"`
	Title    string          `bson:"title"`
	ViewType string          `bson:"view_type"`
	UserName string          `bson:"username"`
	Slug     string          `bson:"slug"`
}
