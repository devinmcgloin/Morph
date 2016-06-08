package schema

import "time"

// Img contains all the proper information for rendering a single photo
type Img struct {
	ID           int       `db:"i_id"`
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
	Album        string    `db:"i_album"`
	CaptureTime  time.Time `db:"i_capture_time"`
	PublishTime  time.Time `db:"i_publish_time"`
	ImgDirection float64   `db:"i_direction"`
	User         int       `db:"u_user"`
	Location     int       `db:"l_id"`
}

type Location struct {
	LID       int     `db:"l_id"`
	Latitude  float64 `db:"l_lat"`
	Longitude float64 `db:"l_lon"`
	Desc      string  `db:"l_desc"`
}

type ImgSource struct {
	ID         int    `db:"s_id"`
	IID        int    `db:"i_id"`
	URL        string `db:"s_url"`
	Resolution int    `db:"s_resolution"`
	Width      int    `db:"s_width"`
	Length     int    `db:"s_length"`
	Size       uint8  `db:"s_size"`
	FileTupe   uint8  `db:"s_file_type"`
}

type User struct {
	ID int `db:"u_id"`
}

// ImgCollection includes a title and collection of Images.
type ImgCollection struct {
	Title  string
	NumImg int
	Images []Img
}
