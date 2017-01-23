package model

import (
	"fmt"

	gj "github.com/sprioc/geojson"
)

// TODO it would be good to have both public and private collections / images.

// Image contains all the proper information for rendering a single photo
type Image struct {
	ShortCode Ref `redis:"-"`

	// Tags         []string `redis:"tags" json:"tags"`
	PublishTime  int64 `redis:"publish_time"`
	LastModified int64 `redis:"last_modified"`
	// Landmarks    []Landmark `redis:"landmarks" json:"landmarks"`
	// Colors       []Color    `redis:"colors" json:"colors"`
	// Labels       []Label    `redis:"labels" json:"labels"`

	Owner     Ref  `redis:"-" json:"owner"`
	Favorites int  `redis:"-" json:"favorites"`
	Featured  bool `redis:"-" json:"featured"`
	Downloads int  `redis:"-" json:"downloads"`
	Views     int  `redis:"-" json:"views"`
	Purchases int  `redis:"-" json:"purchases"`

	// Image Metadata
	Aperture        string    `redis:"aperture" json:"aperture,omitempty"`
	ExposureTime    string    `redis:"exposure_time" json:"exposure_time,omitempty"`
	FocalLength     string    `redis:"focal_length" json:"focal_length,omitempty"`
	ISO             int       `redis:"iso" json:"iso,omitempty"`
	Make            string    `redis:"make" json:"make,omitempty"`
	Model           string    `redis:"model" json:"model,omitempty"`
	LensMake        string    `redis:"lens_make" json:"lens_make,omitempty"`
	LensModel       string    `redis:"lens_model" json:"lens_model,omitempty"`
	PixelXDimension int64     `redis:"pixel_xd" json:"pixel_xd,omitempty"`
	PixelYDimension int64     `redis:"pixel_yd" json:"pixel_yd,omitempty"`
	CaptureTime     int64     `redis:"capture_time" json:"capture_time,omitempty"`
	ImgDirection    float64   `redis:"direction" json:"direction,omitempty"`
	Location        *gj.Point `redis:"-" json:"location,omitempty"`
}

// ImgSource includes the information about the image itself.
type ImgSource struct {
	Thumb  string `redis:"thumb" json:"thumb"`
	Small  string `redis:"small" json:"small"`
	Medium string `redis:"medium" json:"medium"`
	Large  string `redis:"large" json:"large"`
	Raw    string `redis:"raw" json:"raw"`
}

type User struct {
	ShortCode Ref         `redis:"-" json:"username"`
	Email     string      `redis:"email" json:"email"`
	Name      string      `redis:"name" json:"name"`
	Bio       string      `redis:"bio" json:"bio"`
	URL       string      `redis:"personal_site_link" json:"personal_site_link"`
	Location  *gj.Feature `redis:"-" json:"location"`

	Password string `redis:"password" json:"-"`
	Salt     string `redis:"salt" json:"-"`

	Images []Ref `redis:"-" json:"image"`
	// Favorites []Ref `redis:"-" json:"favorited"`

	Featured bool `redis:"-" json:"featured"`
	Admin    bool `redis:"-" json:"admin"`
	Views    int  `redis:"-" json:"views"`

	CreatedAt    int64 `redis:"created_at" json:"created_at"`
	LastModified int64 `redis:"last_modified" json:"last_modified"`

	// Personal Only filled out through /me endpoint
	SeenLinks []Ref `redis:"-" json:"seen_links,omitempty"`
	Purchased []Ref `redis:"-" json:"purchased_links,omitempty"`
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

func (u User) GetTag() string {
	return fmt.Sprintf("%s:%s", "users", u.ShortCode.ShortCode)
}

func (i Image) GetTag() string {
	return fmt.Sprintf("%s:%s", "images", i.ShortCode.ShortCode)
}

func (u User) GetRef() Ref {
	return u.ShortCode
}

func (i Image) GetRef() Ref {
	return i.ShortCode
}
