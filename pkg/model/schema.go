package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	gj "github.com/kpawlik/geojson"
)

// TODO it would be good to have both public and private collections / images.

// TODO need to define external representations of these types and functions
// that map from internal to external.

// Image contains all the proper information for rendering a single photo
type Image struct {
	ID        bson.ObjectId `bson:"_id" json:"-"`
	ShortCode string        `bson:"shortcode" json:"shortcode"`

	MetaData    ImageMetaData `bson:"metadata" json:"metadata"`
	Tags        []string      `bson:"tags" json:"tags"`
	MachineTags []string      `bson:"machine_tags" json:"machine_tags"`
	PublishTime time.Time     `bson:"publish_time" json:"publish_time"`

	Owner       DBRef `bson:"owner" json:"-"`
	OwnerExtern User  `bson:"-" json:"owner"`

	Collections     []DBRef  `bson:"collections" json:"-"`
	CollectionLinks []string `bson:"-" json:"collection_links"`

	FavoritedBy      []DBRef  `bson:"favorited_by" json:"-"`
	FavoritedByLinks []string `bson:"-" json:"favorited_by_links"`

	Sources ImgSource `bson:"sources_link" json:"source_link"`

	Featured  bool `bson:"featured" json:"featured"`
	Downloads int  `bson:"downloads" json:"downloads"`
	Hidden    bool `bson:"hidden" json:"hidden"`
}

type ImageMetaData struct {
	Aperture       Ratio  `bson:"aperture" json:"-"`
	ApertureExtern string `bson:"-" json:"aperture"`

	ExposureTime       Ratio  `bson:"exposure_time" json:"-"`
	ExposureTimeExtern string `bson:"-" json:"exposure_time"`

	FocalLength       Ratio  `bson:"focal_length" json:"-"`
	FocalLengthExtern string `bson:"-" json:"focal_length"`

	ISO             int        `bson:"iso" json:"iso"`
	Orientation     string     `bson:"orientation" json:"orientation"`
	Make            string     `bson:"make" json:"make"`
	Model           string     `bson:"model" json:"model"`
	LensMake        string     `bson:"lens_make" json:"lens_make"`
	LensModel       string     `bson:"lens_model" json:"lens_model"`
	PixelXDimension int64      `bson:"pixel_xd" json:"pixel_xd"`
	PixelYDimension int64      `bson:"pixel_yd" json:"pixel_yd"`
	CaptureTime     time.Time  `bson:"capture_time" json:"capture_time"`
	ImgDirection    float64    `bson:"direction" json:"direction"`
	Location        gj.Feature `bson:"location" json:"location"`
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
	Admin     bool          `bson:"admin" json:"admin"`

	Images     []DBRef  `bson:"images" json:"-"`
	ImageLinks []string `bson:"-" json:"image_links"`

	Collections     []DBRef  `bson:"collections" json:"-"`
	CollectionLinks []string `bson:"-" json:"collection_links"`

	Followes    []DBRef  `bson:"followes" json:"-"`
	FollowLinks []string `bson:"-" json:"follow_links"`

	Favorites     []DBRef  `bson:"favorites" json:"-"`
	FavoriteLinks []string `bson:"-" json:"favorite_links"`

	FollowedBy      []DBRef  `bson:"followed_by" json:"-"`
	FollowedByLinks []string `bson:"-" json:"followed_by_links"`

	FavoritedBy      []DBRef  `bson:"favorited_by" json:"-"`
	FavoritedByLinks []string `bson:"-" json:"favorited_by_links"`

	Email     string     `bson:"email" json:"email"`
	Pass      string     `bson:"password" json:"-"`
	Salt      string     `bson:"salt" json:"-"`
	Name      string     `bson:"name" json:"name"`
	Bio       string     `bson:"bio" json:"bio"`
	URL       string     `bson:"personal_site_link" json:"personal_site_link"`
	Location  gj.Feature `bson:"loc" json:"loc"`
	AvatarURL ImgSource  `bson:"avatar_link" json:"avatar_link"`
}

type Collection struct {
	ID        bson.ObjectId `bson:"_id" json:"-"`
	ShortCode string        `bson:"shortcode" json:"shortcode"`
	Images    []DBRef       `bson:"images" json:"images"`
	Owner     DBRef         `bson:"owner" json:"owner"`

	FollowedBy      []DBRef  `bson:"followed_by" json:"-"`
	FollowedByLinks []string `bson:"-" json:"followed_by_links"`

	FavoritedBy      []DBRef  `bson:"favorited_by" json:"-"`
	FavoritedByLinks []string `bson:"-" json:"favorited_by_links"`

	Desc     string     `bson:"desc" json:"desc"`
	Title    string     `bson:"title" json:"title"`
	ViewType string     `bson:"view_type" json:"view_type"`
	Location gj.Feature `bson:"location" json:"location"`
}

type DBRef struct {
	Collection string `bson:"collection" json:"collection"`
	Shortcode  string `bson:"shortcode" json:"shortcode"`
	Database   string `bson:"db" json:"-"`
}

func (ref DBRef) String() string {
	return getURL(ref)
}
