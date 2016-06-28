package model

import "fmt"

var baseurl = "http://localhost:8080/v0"

func (img *Image) FillExternal(user User) {
	user.FillExternal()
	img.OwnerExtern = user
	img.MetaData.ApertureExtern = img.MetaData.Aperture.Rep
	img.MetaData.FocalLengthExtern = img.MetaData.FocalLength.Rep
	img.MetaData.ExposureTimeExtern = img.MetaData.ExposureTime.Rep

	img.FavoritedByLinks = getURLs(img.FavoritedBy)
	img.CollectionLinks = getURLs(img.Collections)
}

func getURL(ref DBRef) string {
	if ref.Collection != "" && ref.Shortcode != "" {
		return fmt.Sprintf("%s/%s/%s", baseurl, ref.Collection, ref.Shortcode)
	}
	return ""
}

func (usr *User) FillExternal() {
	usr.ImageLinks = getURLs(usr.Images)

	usr.FollowLinks = getURLs(usr.Followes)
	usr.FollowedByLinks = getURLs(usr.FollowedBy)

	usr.FavoriteLinks = getURLs(usr.Favorites)
	usr.FavoritedByLinks = getURLs(usr.FavoritedBy)

	usr.CollectionLinks = getURLs(usr.Collections)
}

func getURLs(refs []DBRef) []string {
	var arr = make([]string, len(refs))
	for i, ref := range refs {
		arr[i] = getURL(ref)
	}
	return arr
}
