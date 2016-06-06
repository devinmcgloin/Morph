package schema

// PhotoMeta contains all the meta information about a specific image
type PhotoMeta struct {
	FStop        int
	ShutterSpeed int
	FOV          int
	ISO          int
}

// Img contains all the proper information for rendering a single photo
type Img struct {
	PID       int
	Title     string
	Desc      string
	URL       string
	Category  string
	PhotoMeta PhotoMeta
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
