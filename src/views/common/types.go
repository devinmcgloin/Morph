package common

import (
	"database/sql"
	"time"
)

///IMAGE VIEW
/// TODO make everything modular and fill out views in the view handlers themselves.
type SingleImg struct {
	IID     uint64         `db:"i_id"`
	Title   sql.NullString `db:"i_title"`
	ImgDesc sql.NullString `db:"i_desc"`
	URL     string         `db:"s_url"`
}

// SingleImgView is a view of a single image on the page.
// ROUTE: /i/:IID
type SingleImgView struct {
	IID          uint64          `db:"i_id"`
	Title        sql.NullString  `db:"i_title"`
	ImgDesc      sql.NullString  `db:"i_desc"`
	Aperture     sql.NullInt64   `db:"i_aperture"`
	ExposureTime sql.NullInt64   `db:"i_exposure_time"`
	FocalLength  sql.NullInt64   `db:"i_focal_length"`
	ISO          sql.NullInt64   `db:"i_iso"`
	CameraBody   sql.NullString  `db:"i_camera_body"`
	Lens         sql.NullString  `db:"i_lens"`
	TagOne       sql.NullString  `db:"i_tag_1"`
	TagTwo       sql.NullString  `db:"i_tag_2"`
	TagThree     sql.NullString  `db:"i_tag_3"`
	CaptureTime  time.Time       `db:"i_capture_time"`
	PublishTime  time.Time       `db:"i_publish_time"`
	Latitude     sql.NullFloat64 `db:"l_lat"`
	Longitude    sql.NullFloat64 `db:"l_lon"`
	LocDesc      sql.NullString  `db:"l_desc"`
	URL          string          `db:"s_url"`
	UID          uint64          `db:"u_id"`
	UserName     string          `db:"u_username"`
	FirstName    sql.NullString  `db:"u_first_name"`
	LastName     sql.NullString  `db:"u_last_name"`
	Bio          sql.NullString  `db:"u_bio"`
	AvatarURL    string          `db:"u_avatar_url"`
}

/// COLLECTION VIEW

// TagCollectionView is a list of images in a collection
// ROUTE: /tag/:tag, /search/*query, /loc/:LID,
type TagCollectionView struct {
	Images []SingleImg
	Tag    string
}

type LocCollectionView struct {
	Images    []SingleImg
	Latitude  sql.NullFloat64 `db:"l_lat"`
	Longitude sql.NullFloat64 `db:"l_lon"`
	Desc      sql.NullString  `db:"l_desc"`
}

type CollectionView struct {
	Images []SingleImg
}

/// USER VIEW

// UserProfileView is a view of a users profile.
// ROUTE: /u/:UID, /settings
type UserProfileView struct {
	UID       uint64          `db:"u_id"`
	UserName  string          `db:"u_username"`
	Email     string          `db:"u_email"`
	FirstName sql.NullString  `db:"u_first_name"`
	LastName  sql.NullString  `db:"u_last_name"`
	Bio       sql.NullString  `db:"u_bio"`
	AvatarURL string          `db:"u_avatar_url"`
	Latitude  sql.NullFloat64 `db:"l_lat"`
	Longitude sql.NullFloat64 `db:"l_lon"`
	Desc      sql.NullString  `db:"l_desc"`

	Images []SingleImg
}

/// ALBUM VIEW

// ROUTE: /album/:AID
type AlbumCollectionView struct {
	Images []SingleImg
}

type UserAuth struct {
}
