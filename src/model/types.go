package model

type DefaultView struct {
	Auth User
}

type SingleImgView struct {
	User  User  `json:"user"`
	Image Image `json:"image"`
	Auth  User  `json:"auth,omitempty"`
}

type TagCollectionView struct {
	Images []Image
	Tags   []string
	Auth   User
}

type LocCollectionView struct {
	Images    []Image
	Locations []Location
	Auth      User
}

type CollectionView struct {
	Images []Image
	Auth   User
}

type UserProfileView struct {
	User   User
	Images []Image
	Auth   User
}

type AlbumCollectionView struct {
	User   User
	Images []Image
	Album  Album
	Auth   User
}
