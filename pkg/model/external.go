package model

import "fmt"

var baseurl = "http://localhost:8080/v0"

func (img *Image) FillExternal(user User) {
	user.FillExternal()
	img.OwnerExtern = user
	img.MetaData.ApertureExtern = img.MetaData.Aperture.Rep
	img.MetaData.FocalLengthExtern = img.MetaData.FocalLength.Rep
	img.MetaData.ExposureTimeExtern = img.MetaData.ExposureTime.Rep
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
	usr.FavoriteLinks = getURLs(usr.Favorites)
}

func getURLs(refs []DBRef) []string {
	var arr []string
	for _, ref := range refs {
		arr = append(arr, getURL(ref))
	}
	return arr
}
