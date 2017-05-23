package model

import (
	"time"

	gj "github.com/sprioc/geojson"
)

// Image contains all the proper information for rendering a single photo
type Image struct {
	Id        int64  `db:"id" json:"-"`
	Shortcode string `db:"shortcode" json:"shortcode"`

	Tags         []string  `db:"-" json:"tags"`
	PublishTime  time.Time `db:"publish_time" json:"publish_time"`
	LastModified time.Time `db:"last_modified" json:"last_modified"`
	// Landmarks    []Landmark `db:"landmarks" json:"landmarks"`
	// Colors       []Color    `db:"colors" json:"colors"`
	//Labels []Label `db:"labels" json:"labels"`

	Owner         int64  `db:"owner_id" json:"-"`
	OwnerUsername string `db:"username" json:"username"`
	Featured      bool   `db:"featured" json:"featured"`
	Downloads     int    `db:"downloads" json:"downloads"`
	Views         int    `db:"views" json:"views"`
	//Purchases int    `db:"" json:"purchases"`
	//Favorites int    `db:"" json:"favorites"`

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
	//ImgDirection    *float64   `db:"direction" json:"direction,omitempty"`
	Location *gj.Point `db:"location" json:"location,omitempty"`
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
	Id       int64   `db:"id" json:"-"`
	Username string  `db:"username" json:"username"`
	Email    string  `db:"email" json:"email"`
	Name     *string `db:"name" json:"name,omitempty"`
	Bio      *string `db:"bio" json:"bio,omitempty"`
	URL      *string `db:"url" json:"url,omitempty"`
	//Location  *gj.Feature `db:"-" json:"location"`

	Password string `db:"password" json:"-"`
	Salt     string `db:"salt" json:"-"`

	Images []string `db:"-" json:"image,omitempty"`
	// Favorites []string `db:"-" json:"favorited"`

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

type Landmark struct {
	Description string
	Location    gj.Point
	Score       float64
}

type Permission string

const (
	CanEdit   = Permission("can_edit")
	CanDelete = Permission("can_delete")
	CanView   = Permission("can_view")
)
