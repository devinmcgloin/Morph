package model

import (
	"time"

	postgis "github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/clr/clr"
)

type User struct {
	Id        int64  `json:"-"`
	Username  string `json:"id"`
	Permalink string `json:"permalink"`

	Email    string          `json:"email"`
	Name     *string         `json:"name,omitempty"`
	Bio      *string         `json:"bio,omitempty"`
	URL      *string         `json:"url,omitempty"`
	Location *postgis.PointS `json:"location,omitempty"`
	Avatars  ImageSource     `json:"avatar_links"`

	ImageLinks    *[]string `json:"images_links,omitempty"`
	FavoriteLinks *[]string `json:"favorite_links,omitempty"`

	FollowedByUser *bool `json:"followed_by_user,omitempty"`

	Featured     bool      `json:"featured"`
	Admin        bool      `json:"admin"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	LastModified time.Time `db:"last_modified" json:"last_modified"`
}

type Image struct {
	Id        int64  `json:"-"`
	Shortcode string `json:"id"`
	Permalink string `json:"permalink"`

	PublishTime  time.Time  `db:"publish_time" json:"publish_time"`
	LastModified time.Time  `db:"last_modified" json:"last_modified"`
	Landmarks    []Landmark `json:"landmarks"`
	Colors       []Color    `json:"colors"`
	Tags         []string   `json:"tags"`
	Labels       []Label    `json:"labels"`

	UserId          int64 `db:"user_id" json:"-"`
	User            *User `json:"user,omitempty"`
	Featured        bool  `json:"featured"`
	FavoritedByUser *bool `json:"favorited_by_user,omitempty"`

	Stats    ImageStats    `json:"stats"`
	Source   ImageSource   `json:"src_links"`
	Metadata ImageMetadata `json:"metadata"`
}

type ImageStats struct {
	Downloads int `json:"downloads"`
	Views     int `json:"views"`
	Favorites int `json:"favorites"`
}

type ImageSource struct {
	Thumb  string `json:"thumb"`
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
	Raw    string `json:"raw"`
}

type ImageMetadata struct {
	Aperture        *string         `db:"aperture" json:"aperture,omitempty"`
	ExposureTime    *string         `db:"exposure_time" json:"exposure_time,omitempty"`
	FocalLength     *string         `db:"focal_length" json:"focal_length,omitempty"`
	ISO             *int            `db:"iso" json:"iso,omitempty"`
	Make            *string         `db:"make" json:"make,omitempty"`
	Model           *string         `db:"model" json:"model,omitempty"`
	LensMake        *string         `db:"lens_make" json:"lens_make,omitempty"`
	LensModel       *string         `db:"lens_model" json:"lens_model,omitempty"`
	PixelXDimension *int64          `db:"pixel_xd" json:"pixel_xd,omitempty"`
	PixelYDimension *int64          `db:"pixel_yd" json:"pixel_yd,omitempty"`
	CaptureTime     *time.Time      `db:"capture_time" json:"capture_time,omitempty"`
	ImageDirection  *float64        `db:"dir" json:"direction,omitempty"`
	Location        *postgis.PointS `db:"loc" json:"location,omitempty"`
}

type Color struct {
	SRGB          clr.RGB `json:"sRGB"`
	Hex           string
	HSV           clr.HSV
	Shade         clr.ColorSpace `db:"shade" json:"shade"`
	ColorName     clr.ColorSpace `db:"color" json:"color_name"`
	PixelFraction float64        `json:"pixel_fraction"`
	Score         float64        `json:"score"`
}

type Label struct {
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}

type Landmark struct {
	Description string         `json:"description"`
	Location    postgis.PointS `json:"location"`
	Score       float64        `json:"score"`
}
