package model

import (
	"time"

	gj "github.com/kpawlik/geojson"
	"gopkg.in/mgo.v2/bson"
)

// TODO it would be good to have both public and private collections / images.

// TODO need to define external representations of these types and functions
// that map from internal to external.

// Image contains all the proper information for rendering a single photo
type Image struct {
	ID        bson.ObjectId `bson:"_id" json:"_id"`
	ShortCode string        `bson:"shortcode" json:"shortcode"`

	MetaData    ImageMetaData `bson:"metadata"`
	Tags        []string      `bson:"tags,omitempty" json:"tags,omitempty"`
	MachineTags []string      `bson:"machine_tags,omitempty" json:"machine_tags,omitempty"`
	PublishTime time.Time     `bson:"publish_time" json:"publish_time"`

	User DBRef `bson:"user" json:"user"`

	AlbumID      DBRef   `bson:"album_id,omitempty" json:"album_id,omitempty"`
	EventID      []DBRef `bson:"event_id,omitempty" json:"event_id,omitempty"`
	CollectionID []DBRef `bson:"collection_id,omitempty" json:"collection_id,omitempty"`

	Sources ImgSource `bson:"sources" json:"sources"`

	Featured  bool    `bson:"featured" json:"featured"`
	Downloads int     `bson:"downloads"`
	Favorites []DBRef `bson:"favoriters"`
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
type ImgSource struct {
	Thumb  string `bson:"thumb" json:"thumb"`
	Small  string `bson:"small" json:"small"`
	Medium string `bson:"medium" json:"medium"`
	Large  string `bson:"large" json:"large"`
	Raw    string `bson:"raw" json:"raw"`
}

type User struct {
	ID         bson.ObjectId `bson:"_id" json:"-"`
	ShortCode  string        `bson:"shortcode" json:"shortcode"`
	Admin      bool          `bson:"admin" json:"admin"`
	Images     []DBRef       `bson:"images" json:"images,omitempty"`
	Followes   []DBRef       `bson:"followes" json:"followes,omitempty"`
	Favorites  []DBRef       `bson:"favorites" json:"favorites,omitempty"`
	Email      string        `bson:"email" json:"email"`
	Pass       string        `bson:"password" json:"-"`
	Salt       string        `bson:"salt" json:"-"`
	Name       string        `bson:"name" json:"name,omitempty"`
	Bio        string        `bson:"bio,omitempty" json:"bio,omitempty"`
	URL        string        `bson:"url" json:"url,omitempty"`
	Location   gj.Feature    `bson:"loc" json:"loc,omitempty"`
	AvatarURL  ImgSource     `bson:"avatar_url" json:"avatar_url"`
	Followers  []DBRef       `bson:"followers" json:"followers,omitempty"`
	Favoriters []DBRef       `bson:"favoriters" json:"favoriters,omitempty"`
}

// Albums need to support custom ordering, have to intgrate this on a per image
// level, or create sperate field for index.
type Album struct {
	ID        bson.ObjectId `bson:"_id" json:"-"`
	ShortCode string        `bson:"shortcode" json:"shortcode"`
	Images    []DBRef       `bson:"images" json:"images,omitempty"`
	Desc      string        `bson:"desc,omitempty" json:"desc,omitempty"`
	Title     string        `bson:"title" json:"title"`
	ViewType  string        `bson:"view_type" json:"view_type"`
	User      DBRef         `bson:"user" json:"user"`
	Followers []DBRef       `bson:"followers" json:"followers,omitempty"`
	Favorites []DBRef       `bson:"favoriters" json:"favoriters,omitempty"`
}

type Event struct {
	ID        bson.ObjectId `bson:"_id" json:"-"`
	ShortCode string        `bson:"shortcode" json:"shortcode"`
	Images    []DBRef       `bson:"images" json:"images,omitempty"`
	Desc      string        `bson:"desc,omitempty" json:"desc,omitempty"`
	Title     string        `bson:"title" json:"title"`
	ViewType  string        `bson:"view_type" json:"view_type"`
	Followers []DBRef       `bson:"followers" json:"followers,omitempty"`
	Favorites []DBRef       `bson:"favoriters" json:"favoriters,omitempty"`
	Location  gj.Feature    `bson:"location" json:"location"`
	TimeStart time.Time     `bson:"time_start" json:"time_start"`
	TimeEnd   time.Time     `bson:"time_end" json:"time_end"`
}

type Collection struct {
	ID        bson.ObjectId `bson:"_id" json:"-"`
	ShortCode string        `bson:"shortcode" json:"shortcode"`
	Images    []DBRef       `bson:"images" json:"images,omitempty"`
	Curator   DBRef         `bson:"curator" json:"curator"`
	Users     []DBRef       `bson:"users" json:"users"`
	Desc      string        `bson:"desc,omitempty" json:"desc,omitempty"`
	Title     string        `bson:"title" json:"title"`
	ViewType  string        `bson:"view_type" json:"view_type"`
	Followers []DBRef       `bson:"followers" json:"followers,omitempty"`
	Favorites []DBRef       `bson:"favoriters" json:"favoriters,omitempty"`
}

type DBRef struct {
	Collection string `bson:"collection" json:"collection"`
	Shortcode  string `bson:"shortcode" json:"shortcode"`
	Database   string `bson:"db,omitempty" json:"-"`
}
