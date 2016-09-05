package model

import (
	"googlemaps.github.io/maps"

	"gopkg.in/mgo.v2/bson"

	gj "github.com/sprioc/geojson"
)

// TODO it would be good to have both public and private collections / images.

// Image contains all the proper information for rendering a single photo
type Image struct {
	ID        bson.ObjectId `bson:"_id" json:"-"`
	ShortCode string        `bson:"shortcode" json:"shortcode"`

	MetaData    ImageMetaData `bson:"metadata" json:"metadata"`
	Tags        []string      `bson:"tags" json:"tags"`
	PublishTime int64         `bson:"publish_time" json:"publish_time"`
	Sources     ImgSource     `bson:"sources_link" json:"source_link"`
	Landmarks   []Landmark    `bson:"landmarks" json:"landmarks"`
	Colors      []Color       `bson:"colors" json:"colors"`
	Labels      []Label       `bson:"labels" json:"labels"`

	OwnerLink         string   `bson:"-" json:"owner"`
	CollectionInLinks []string `bson:"-" json:"collection_links"`
	FavoritedByLinks  []string `bson:"-" json:"favorited_by_links"`
	Featured          bool     `bson:"-" json:"featured"`
	Downloads         int      `bson:"-" json:"downloads"`
	Views             int      `bson:"-" json:"views"`
	Purchases         int      `bson:"-" json:"purchases"`
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
	ShortCode string        `bson:"shortcode" json:"shortcode"`
	Email     string        `bson:"email" json:"email"`
	Name      string        `bson:"name" json:"name"`
	Bio       string        `bson:"bio" json:"bio,omitempty"`
	URL       string        `bson:"personal_site_link" json:"personal_site_link,omitempty"`
	Location  *gj.Feature   `bson:"location" json:"location,omitempty"`
	AvatarURL ImgSource     `bson:"avatar_link" json:"avatar_link"`

	ImageLinks      []string `bson:"-" json:"iamge_links"`
	CollectionLinks []string `bson:"-" json:"collection_links"`

	FollowedLinks  []string `bson:"-" json:"followed_links"`
	FavoritedLinks []string `bson:"-" json:"favorited_links"`

	Featured bool `bson:"-" json:"featured"`
	Admin    bool `bson:"-" json:"admin"`
	Views    int  `bson:"-" json:"views"`

	// Personal Only filled out through /me endpoint
	FavoritedByLinks []string `bson:"-" json:"favorited_by_links,omitempty"`
	FollowedByLinks  []string `bson:"-" json:"followed_by_links,omitempty"`

	SeenLinks []string `bson:"-" json:"seen_links,omitempty"`
	Purchased []string `bson:"-" json:"purchased_links,omitempty"`
}

type Collection struct {
	ID        bson.ObjectId `bson:"_id" json:"-"`
	ShortCode string        `bson:"shortcode" json:"shortcode"`
	Desc      string        `bson:"desc" json:"desc,omitempty"`
	Title     string        `bson:"title" json:"title,omitempty"`
	Location  *gj.Feature   `bson:"location" json:"location,omitempty"`

	OwnerLink        string   `bson:"-" json:"owner"`
	ImageLinks       []string `bson:"-" json:"iamge_links"`
	FavoritedByLinks []string `bson:"-" json:"favorited_by_links"`
	FollowedByLinks  []string `bson:"-" json:"followed_by_links"`
	Featured         bool     `bson:"-" json:"featured"`
	Views            int      `bson:"-" json:"views"`
	ViewType         string   `bson:"-" json:"view_type"`
}

type Location struct {
	GoogleLoc maps.GeocodingResult `bson:"google_geo"`
	Bounds    gj.Polygon           `bson:"bounds"`
}

type Landmark struct {
	Description string
	Location    gj.Point
	Score       float64
}

type Color struct {
	Color struct {
		Red   int
		Green int
		Blue  int
	}
	PixelFraction float64
	Score         float64
}

type Label struct {
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}
