package SQL

// TODO consider what can be null and what cannot and set DB settings to reflect
// this.

// sql.Nullxxx Have a valid field and a value field. Drive knows how to handle
// them. Just need to check before working with these types.

import (
	"database/sql"
	"strconv"
	"time"
)

// Img contains all the proper information for rendering a single photo
type Img struct {
	IID          uint64          `db:"i_id"`
	Title        sql.NullString  `db:"i_title"`
	Desc         sql.NullString  `db:"i_desc"`
	Aperture     sql.NullInt64   `db:"i_aperture"`
	ExposureTime sql.NullInt64   `db:"i_exposure_time"`
	FocalLength  sql.NullInt64   `db:"i_focal_length"`
	ISO          sql.NullInt64   `db:"i_iso"`
	Orientation  sql.NullString  `db:"i_orientation"`
	CameraBody   sql.NullString  `db:"i_camera_body"`
	Lens         sql.NullString  `db:"i_lens"`
	TagOne       sql.NullString  `db:"i_tag_1"`
	TagTwo       sql.NullString  `db:"i_tag_2"`
	TagThree     sql.NullString  `db:"i_tag_3"`
	AID          sql.NullInt64   `db:"a_id"`
	CaptureTime  time.Time       `db:"i_capture_time"`
	PublishTime  time.Time       `db:"i_publish_time"`
	ImgDirection sql.NullFloat64 `db:"i_direction"`
	UID          sql.NullInt64   `db:"u_id"`
	LID          sql.NullInt64   `db:"l_id"`
}

type Location struct {
	LID       uint64          `db:"l_id"`
	Latitude  sql.NullFloat64 `db:"l_lat"`
	Longitude sql.NullFloat64 `db:"l_lon"`
	Desc      sql.NullString  `db:"l_desc"`
}

// ImgSource includes the information about the image itself.
// Size indicates how large the image is.
type ImgSource struct {
	SID        uint64        `db:"s_id"`
	IID        uint64        `db:"i_id"`
	URL        string        `db:"s_url"`
	Resolution sql.NullInt64 `db:"s_resolution"`
	Width      sql.NullInt64 `db:"s_width"`
	Height     sql.NullInt64 `db:"s_height"`
	Size       string        `db:"s_size"`
	FileType   string        `db:"s_file_type"`
}

type User struct {
	UID       uint64         `db:"u_id"`
	Usernmae  string         `db:"u_username"`
	Email     string         `db:"u_email"`
	FirstName string         `db:"u_first_name"`
	LastName  string         `db:"u_last_name"`
	Bio       sql.NullString `db:"u_bio"`
	LID       sql.NullInt64  `db:"u_loc"`
	AvatarURL string         `db:"u_avatar_url"`
}

type Album struct {
	AID      uint64         `db:"a_id"`
	UID      uint64         `db:"u_id"`
	Desc     sql.NullString `db:"a_desc"`
	Title    string         `db:"a_title"`
	ViewType string         `db:"a_view_type"`
}

///IMAGE VIEW

// FeatureSingleImgView is a view of a single image on the page.
// ROUTE: /i/:ID
type SingleImgView struct {
	User
	Img
	ImgSource
}

/// COLLECTION VIEW

// CollectionView is a list of images in a collection
// ROUTE: /tag/:tag, /search/*query, /loc/:LID, /tag/:tag/:IID
type CollectionView struct {
	Images []SingleImgView
}

/// USER VIEW

// UserProfileView is a view of a users profile.
// ROUTE: /u/:UID, /settings
type UserProfileView struct {
	User
	Location
	Images []SingleImgView
}

/// ALBUM VIEW

// ROUTE: /album/:AID
type AlbumCollectionView struct {
	User
	Images []SingleImgView
	Album
}

// type Image struct {
// 	IID          int
// 	UID          int
// 	LID          int
// 	URL          string
// 	Title        string
// 	Description  string
// 	CaptureTime  time.Time
// 	Publishtime  time.Time
// 	Tag1         string
// 	Tag2         string
// 	Tag3         string
// 	Lon          float64
// 	Lat          float64
// 	LocationDesc string
// }
//
// type User struct {
// 	UID       int
// 	LID       int
// 	FirstName string
// 	LastName  string
// 	Bio       string
// 	Location  string
// 	AvatarURL string
// }
//
// type Album struct {
// 	AID      int
// 	Desc     string
// 	Title    string
// 	ViewType string
// }

func ToNullInt64(s string) sql.NullInt64 {
	i, err := strconv.Atoi(s)
	return sql.NullInt64{Int64: int64(i), Valid: err == nil}
}

func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
