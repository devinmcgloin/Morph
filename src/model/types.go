package model

type SingleImgView struct {
	User  User
	Image Image
}

type TagCollectionView struct {
	Images []Image
	Tags   []string
}

type LocCollectionView struct {
	Images    []Image
	Locations []Location
}

type CollectionView struct {
	Images []Image
}

type UserProfileView struct {
	User   User
	Images []Image
}

type AlbumCollectionView struct {
	User   User
	Images []Image
	Album  Album
}

type UserAuth struct {
}
