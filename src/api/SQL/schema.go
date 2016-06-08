package SQL

import "time"

// Img contains all the proper information for rendering a single photo
type Img struct {
	IID          int       `db:"i_id"`
	Title        string    `db:"i_title"`
	Desc         string    `db:"i_desc"`
	Aperture     int       `db:"i_aperture"`
	ExposureTime int       `db:"i_exposure_time"`
	FocalLength  int       `db:"i_focal_length"`
	ISO          int       `db:"i_iso"`
	Orientation  string    `db:"i_orientation"`
	CameraBody   string    `db:"i_camera_body"`
	Lens         string    `db:"i_lens"`
	TagOne       string    `db:"i_tag_1"`
	TagTwo       string    `db:"i_tag_2"`
	TagThree     string    `db:"i_tag_3"`
	AID          int       `db:"a_id"`
	CaptureTime  time.Time `db:"i_capture_time"`
	PublishTime  time.Time `db:"i_publish_time"`
	ImgDirection float64   `db:"i_direction"`
	UID          int       `db:"u_id"`
	LID          int       `db:"l_id"`
}

type Location struct {
	LID       int     `db:"l_id"`
	Latitude  float64 `db:"l_lat"`
	Longitude float64 `db:"l_lon"`
	Desc      string  `db:"l_desc"`
}

// ImgSource includes the information about the image itself.
// Size indicates how large the image is.
type ImgSource struct {
	SID        int    `db:"s_id"`
	IID        int    `db:"i_id"`
	URL        string `db:"s_url"`
	Resolution int    `db:"s_resolution"`
	Width      int    `db:"s_width"`
	Height     int    `db:"s_height"`
	Size       string `db:"s_size"`
	FileType   string `db:"s_file_type"`
}

type User struct {
	UID       int    `db:"u_id"`
	Usernmae  string `db:"u_username"`
	Email     string `db:"u_email"`
	FirstName string `db:"u_first_name"`
	LastName  string `db:"u_first_name"`
	Bio       string `db:"u_bio"`
	LID       int    `db:"u_loc"`
	AvatarURL string `db:"u_avatar_url"`
}

type Album struct {
	AID      int    `db:"a_id"`
	UID      int    `db:"u_id"`
	Desc     string `db:"a_desc"`
	Title    string `db:"a_title"`
	ViewType string `db:"a_view_type"`
}
