package refs

import (
	"fmt"
	"log"
	"strings"

	"github.com/sprioc/sprioc-core/pkg/env"
	"github.com/sprioc/sprioc-core/pkg/model"
)

var dbname = env.Getenv("MONGODB_NAME", "sprioc")
var baseurl = "http://localhost:8080/api/v0/"

func GetImageRef(shortcode string) model.DBRef {
	return model.DBRef{Database: dbname, Collection: "images", Shortcode: shortcode}
}
func GetCollectionRef(shortcode string) model.DBRef {
	return model.DBRef{Database: dbname, Collection: "collections", Shortcode: shortcode}
}
func GetAlbumRef(shortcode string) model.DBRef {
	return model.DBRef{Database: dbname, Collection: "albums", Shortcode: shortcode}
}
func GetUserRef(shortcode string) model.DBRef {
	return model.DBRef{Database: dbname, Collection: "users", Shortcode: shortcode}
}

func GetURL(ref model.DBRef) string {
	if ref.Collection != "" && ref.Shortcode != "" {
		return fmt.Sprintf("%s%s/%s", baseurl, ref.Collection, ref.Shortcode)
	}
	return ""
}

func GetURLs(refs []model.DBRef) []string {
	var arr = make([]string, len(refs))
	for i, ref := range refs {
		arr[i] = GetURL(ref)
	}
	return arr
}

func GetRef(url string) model.DBRef {

	parts := strings.Split(strings.Replace(url, baseurl, "", 1), "/")
	log.Println(parts, len(parts))
	return model.DBRef{Database: dbname, Collection: parts[0], Shortcode: parts[1]}
}

func GetRefs(urls []string) []model.DBRef {
	var arr = make([]model.DBRef, len(urls))
	for i, ref := range urls {
		arr[i] = GetRef(ref)
	}
	return arr
}

func FillExternalImage(img *model.Image) {
	img.OwnerLink = GetURL(img.Owner)
	img.MetaData.ApertureExtern = img.MetaData.Aperture.Rep
	img.MetaData.FocalLengthExtern = img.MetaData.FocalLength.Rep
	img.MetaData.ExposureTimeExtern = img.MetaData.ExposureTime.Rep

	img.FavoritedByLinks = GetURLs(img.FavoritedBy)
	img.CollectionLinks = GetURLs(img.Collections)
}

func FillExternalUser(usr *model.User) {
	usr.ImageLinks = GetURLs(usr.Images)

	usr.FollowLinks = GetURLs(usr.Followes)
	usr.FollowedByLinks = GetURLs(usr.FollowedBy)

	usr.FavoriteLinks = GetURLs(usr.Favorites)
	usr.FavoritedByLinks = GetURLs(usr.FavoritedBy)

	usr.CollectionLinks = GetURLs(usr.Collections)
}

func FillExternalCollection(col *model.Collection, user model.User) {
	FillExternalUser(&user)

	col.OwnerExtern = user
	col.ImageLinks = GetURLs(col.Images)
	col.FavoritedByLinks = GetURLs(col.FavoritedBy)
	col.FollowedByLinks = GetURLs(col.FollowedBy)
}
