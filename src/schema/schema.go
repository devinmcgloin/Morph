package schema

import "time"

// Img contains all the proper information for rendering a single photo
type Img struct {
	IID          int       `db:"i_id"`
	Title        string    `db:"i_title"`
	Desc         string    `db:"i_desc"`
	URL          string    `db:"i_url"`
	FStop        int       `db:"i_fstop"`
	ShutterSpeed int       `db:"i_shutter_speed"`
	FOV          int       `db:"i_fov"`
	ISO          int       `db:"i_iso"`
	Category     string    `db:"i_category"`
	PublishDate  time.Time `db:"i_publish_date"`
}

// ImgCollection includes a title and collection of Images.
type ImgCollection struct {
	Title  string
	NumImg int
	Images []Img
}

// PageMeta is a type for the meta tags found at the top of the page.
type PageMeta struct {
	Title         string
	Image         string
	URL           string
	Description   string
	Author        string
	AuthorTwitter string
	Keywords      []string
}
