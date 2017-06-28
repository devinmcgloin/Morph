package model

import (
	"time"

	"github.com/sprioc/clr/clr"
	gj "github.com/sprioc/geojson"
)

// Image contains all the proper information for rendering a single photo
type Image struct {
	Id        int64  `db:"id" json:"-"`
	Shortcode string `db:"shortcode" json:"permalink"`

	Tags         []string   `db:"-" json:"tags"`
	PublishTime  time.Time  `db:"publish_time" json:"publish_time"`
	LastModified time.Time  `db:"last_modified" json:"last_modified"`
	Landmarks    []Landmark `db:"-" json:"landmarks"`
	Colors       []Color    `db:"-" json:"colors"`
	Labels       []Label    `db:"-" json:"labels"`

	Owner         int64  `db:"owner_id" json:"-"`
	OwnerUsername string `db:"username" json:"owner_link"`
	Featured      bool   `db:"featured" json:"featured"`
	Downloads     int    `db:"downloads" json:"downloads"`
	Views         int    `db:"views" json:"views"`
	Favorites     int    `db:"" json:"favorites"`

	// Image Metadata
	Aperture        *string    `db:"aperture" json:"aperture,omitempty"`
	ExposureTime    *string    `db:"exposure_time" json:"exposure_time,omitempty"`
	FocalLength     *string    `db:"focal_length" json:"focal_length,omitempty"`
	ISO             *int       `db:"iso" json:"iso,omitempty"`
	Make            *string    `db:"make" json:"make,omitempty"`
	Model           *string    `db:"model" json:"model,omitempty"`
	LensMake        *string    `db:"lens_make" json:"lens_make,omitempty"`
	LensModel       *string    `db:"lens_model" json:"lens_model,omitempty"`
	PixelXDimension *int64     `db:"pixel_xd" json:"pixel_xd,omitempty"`
	PixelYDimension *int64     `db:"pixel_yd" json:"pixel_yd,omitempty"`
	CaptureTime     *time.Time `db:"capture_time" json:"capture_time,omitempty"`
	ImgDirection    *float64   `db:"direction" json:"direction,omitempty"`
	Location        *gj.Point  `db:"location" json:"location,omitempty"`
}

func (i Image) GetRef() Ref {
	return Ref{Collection: Images, Id: i.Id, Shortcode: i.Shortcode}
}

// ImgSource includes the information about the image itself.
type ImgSource struct {
	Thumb  string `json:"thumb"`
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
	Raw    string `json:"raw"`
}

type User struct {
	Id       int64       `db:"id" json:"-"`
	Username string      `db:"username" json:"permalink"`
	Email    string      `db:"email" json:"email"`
	Name     *string     `db:"name" json:"name,omitempty"`
	Bio      *string     `db:"bio" json:"bio,omitempty"`
	URL      *string     `db:"url" json:"url,omitempty"`
	Location *gj.Feature `db:"-" json:"location"`

	Password string `db:"password" json:"-"`
	Salt     string `db:"salt" json:"-"`

	Images    []string `db:"-" json:"image_links"`
	Favorites []string `db:"-" json:"favorited_links"`

	Featured     bool      `db:"featured" json:"featured"`
	Admin        bool      `db:"admin" json:"admin"`
	Views        int       `db:"views" json:"views"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	LastModified time.Time `db:"last_modified" json:"last_modified"`
}

func (u User) GetRef() Ref {
	return Ref{Collection: Users, Id: u.Id, Shortcode: u.Username}
}

type Color struct {
	Color         clr.RGB `json:"color"`
	PixelFraction float64 `json:"pixel_fraction"`
	Score         float64 `json:"score"`
}

type Label struct {
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}

type Landmark struct {
	Description string   `json:"description"`
	Location    gj.Point `json:"location"`
	Score       float64  `json:"score"`
}

type Permission string

const (
	CanEdit   = Permission("can_edit")
	CanDelete = Permission("can_delete")
	CanView   = Permission("can_view")
)
