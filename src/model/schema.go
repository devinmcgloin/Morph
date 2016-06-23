package model

import (
	"time"

	gj "github.com/kpawlik/geojson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Image contains all the proper information for rendering a single photo
type Image struct {
	ID bson.ObjectId `bson:"_id" json:"_id"`

	MetaData    ImageMetaData `bson:"metadata"`
	Tags        []string      `bson:"tags,omitempty" json:"tags,omitempty"`
	MachineTags []string      `bson:"machine_tags,omitempty" json:"machine_tags,omitempty"`
	PublishTime time.Time     `bson:"publish_time" json:"publish_time"`

	User mgo.DBRef `bson:"user" json:"user"`

	AlbumID      mgo.DBRef   `bson:"album_id,omitempty" json:"album_id,omitempty"`
	EventID      []mgo.DBRef `bson:"event_id,omitempty" json:"event_id,omitempty"`
	CollectionID []mgo.DBRef `bson:"collection_id,omitempty" json:"collection_id,omitempty"`

	Sources ImgSource `bson:"sources" json:"sources"`

	Featured  bool        `bson:"featured,omitempty" json:"featured,omitempty"`
	Downloads int         `bson:"downloads"`
	Favorites []mgo.DBRef `bson:"favoriters"`
}

type ImageMetaData struct {
	Aperture        Ratio      `bson:"aperture,omitempty" json:"aperture,omitempty"`
	ExposureTime    Ratio      `bson:"exposure_time,omitempty" json:"exposure_time,omitempty"`
	FocalLength     Ratio      `bson:"focal_length,omitempty" json:"focal_length,omitempty"`
	ISO             int        `bson:"iso,omitempty" json:"iso,omitempty"`
	Orientation     string     `bson:"orientation,omitempty" json:"orientation,omitempty"`
	Make            string     `bson:"make,omitempty" json:"make,omitempty"`
	Model           string     `bson:"model,omitempty" json:"model,omitempty"`
	LensMake        string     `bson:"lens_make,omitempty" json:"lens_make,omitempty"`
	LensModel       string     `bson:"lens_model,omitempty" json:"lens_model,omitempty"`
	PixelXDimension int64      `bson:"pixel_xd,omitempty" json:"pixel_xd,omitempty"`
	PixelYDimension int64      `bson:"pixel_yd,omitempty" json:"pixel_yd,omitempty"`
	CaptureTime     time.Time  `bson:"capture_time" json:"capture_time"`
	ImgDirection    float64    `bson:"direction,omitempty" json:"direction,omitempty"`
	Location        gj.Feature `bson:"location,omitempty" json:"location,omitempty"`
}

// ImgSource includes the information about the image itself.
// Size indicates how large the image is.
type ImgSource struct {
	Thumb  URL
	Small  URL
	Medium URL
	Large  URL
	Raw    URL
}

type User struct {
	ID         bson.ObjectId `bson:"_id"`
	Images     []mgo.DBRef   `bson:"images"`
	Followes   []mgo.DBRef   `bson:"followes"`
	Favorites  []mgo.DBRef   `bson:"favorites"`
	UserName   UserName      `bson:"username"`
	Email      string        `bson:"email"`
	Pass       string        `bson:"-" json:"-"`
	Salt       string        `bson:"-" json:"-"`
	Name       string        `bson:"name"`
	Bio        string        `bson:"bio,omitempty"`
	URL        URL           `bson:"url"`
	Location   gj.Feature    `bson:"loc"`
	AvatarURL  ImgSource     `bson:"avatar_url"`
	Provider   string        `bson:"provider"`
	Followers  []mgo.DBRef   `bson:"followers"`
	Favoriters []mgo.DBRef   `bson:"favoriters"`
}

// Albums need to support custom ordering, have to intgrate this on a per image
// level, or create sperate field for index.
type Album struct {
	ID        bson.ObjectId `bson:"_id"`
	Images    []mgo.DBRef   `bson:"images"`
	Desc      string        `bson:"desc"`
	Title     string        `bson:"title"`
	ViewType  string        `bson:"view_type"`
	User      mgo.DBRef     `bson:"username"`
	Followers []mgo.DBRef   `bson:"followers"`
	Favorites []mgo.DBRef   `bson:"favorites"`
}

type Event struct {
	ID        bson.ObjectId `bson:"_id"`
	Images    []mgo.DBRef   `bson:"images"`
	Desc      string        `bson:"desc"`
	Title     string        `bson:"title"`
	ViewType  string        `bson:"view_type"`
	UserName  mgo.DBRef     `bson:"username"`
	Followers []mgo.DBRef   `bson:"followers"`
	Favorites []mgo.DBRef   `bson:"favorites"`
	Location  gj.Feature    `bson:"location"`
	TimeStart time.Time     `bson:"time_start"`
	TimeEnd   time.Time     `bson:"time_end"`
}

type Collection struct {
	ID        bson.ObjectId `bson:"_id"`
	Images    []mgo.DBRef   `bson:"images"`
	Users     []mgo.DBRef   `bson:"users"`
	Desc      string        `bson:"desc"`
	Title     string        `bson:"title"`
	ViewType  string        `bson:"view_type"`
	Followers []mgo.DBRef   `bson:"followers"`
	Favorites []mgo.DBRef   `bson:"favorites"`
}
