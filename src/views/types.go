package views

import "time"

/// BUILDING BLOCK TYPES

type Image struct {
	IID          int
	UID          int
	LID          int
	URL          string
	Title        string
	Description  string
	CaptureTime  time.Time
	Publishtime  time.Time
	Tag1         string
	Tag2         string
	Tag3         string
	Lon          float64
	Lat          float64
	LocationDesc string
}

type User struct {
	UID       int
	LID       int
	FirstName string
	LastName  string
	Bio       string
	Location  string
	AvatarURL string
}

///IMAGE VIEW

// FeatureSingleImgView is a view of a single image on the page.
type FeatureSingleImgView struct {
	User  User
	Image Image
}

/// COLLECTION VIEW

// CollectionView is a list of images in a collection
type CollectionView struct {
	Images []Image
}

// CollectionFeatureView in which one photo has been selected but others are still
// present in background.
type CollectionFeatureView struct {
	focusIndex int
	Images     []Image
}

/// USER VIEW

// UserProfileView is a view of a users profile.
type UserProfileView struct {
	User   User
	Images []Image
}

/// ALBUM VIEW

type AlbumCollectionView struct {
	Author User
	Images []Image
}
