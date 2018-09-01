package image

import (
	"context"
	"time"

	postgis "github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/clr/clr"
)

//go:generate moq -out image_service_runner.go . ImageService

type Service interface {
	ImageByID(ctx context.Context, id uint64) (*Image, error)
	ImageByShortcode(ctx context.Context, shortcode string) (*Image, error)
	ExistsByShortcode(ctx context.Context, shortcode string) (bool, error)
	CreateImage(ctx context.Context, i *Image) error
	DeleteImage(ctx context.Context, id uint64) error

	RandomImage(ctx context.Context) (*Image, error)
	RandomImageForUser(ctx context.Context, userID uint64) (*Image, error)

	Feature(ctx context.Context, id uint64, user uint64) error
	UnFeature(ctx context.Context, id uint64, user uint64) error

	ImagesForUser(ctx context.Context, id uint64) (*[]Image, error)

	ImageStats(ctx context.Context, id uint64) (*ImageStats, error)
	ImageSource(ctx context.Context, id uint64) (*ImageSource, error)
	ImageMetadata(ctx context.Context, id uint64) (*ImageMetadata, error)
	SetImageMetadata(ctx context.Context, id uint64, metadata ImageMetadata) error
	ImageLocation(ctx context.Context, id uint64) (*Location, error)
	SetImageLocation(ctx context.Context, id uint64, location Location) error
	ImageColors(ctx context.Context, id uint64) (*[]Color, error)
	SetImageColors(ctx context.Context, id uint64, colors []Color) error
	ImageTags(ctx context.Context, id uint64) (*[]string, error)
	ImageLabels(ctx context.Context, id uint64) (*[]Label, error)
}

type Image struct {
	ID        uint64
	Shortcode string
	UserID    uint64

	Featured bool

	PublishedAt  time.Time
	LastModified time.Time

	Title       *string
	Description *string
}

type ImageStats struct {
	Views     uint64
	Downloads uint64
	Favorites uint64
}

type ImageSource struct {
	Thumb  string `json:"thumb"`
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
	Raw    string `json:"raw"`
}

type ImageMetadata struct {
	Aperture        *float64   `db:"aperture" json:"aperture,omitempty"`
	ExposureTime    *float64   `db:"exposure_time" json:"exposure_time,omitempty"`
	FocalLength     *int       `db:"focal_length" json:"focal_length,omitempty"`
	ISO             *int       `db:"iso" json:"iso,omitempty"`
	Make            *string    `db:"make" json:"make,omitempty"`
	Model           *string    `db:"model" json:"model,omitempty"`
	LensMake        *string    `db:"lens_make" json:"lens_make,omitempty"`
	LensModel       *string    `db:"lens_model" json:"lens_model,omitempty"`
	PixelXDimension uint64     `db:"pixel_xd" json:"pixel_xd"`
	PixelYDimension uint64     `db:"pixel_yd" json:"pixel_yd"`
	CaptureTime     *time.Time `db:"capture_time" json:"capture_time,omitempty"`
	Location        *Location  `db:"location" json:"location,omitempty"`
	Orientation     uint16     `db:"-" json:"-"`
}

type Location struct {
	ImageDirection *float64        `db:"dir" json:"direction,omitempty"`
	Point          *postgis.PointS `db:"loc" json:"-"`
	LatLng         *Point          `json:"point,omitempty"`
	Description    *string         `db:"description" json:"description"`
}

type Color struct {
	SRGB          clr.RGB        `json:"sRGB"`
	Hex           string         `json:"hex"`
	HSV           clr.HSV        `json:"hsv"`
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

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
