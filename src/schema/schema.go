package schema

// PhotoMeta contains all the meta information about a specific image
type PhotoMeta struct {
	FStop        int
	ShutterSpeed int
	FOV          int
	ISO          int
}

// ImgPage contains all the proper information for rendering a single photo
type ImgPage struct {
	Title     string
	Desc      string
	ImgURL    string
	PhotoMeta PhotoMeta
	PageMeta  PageMeta
}

// ImgCollection includes a title and collection of Images.
type ImgCollection struct {
	Title    string
	NumImg   int
	Images   []ImgPage
	PageMeta PageMeta
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
