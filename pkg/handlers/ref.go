package handlers

import (
	"github.com/sprioc/sprioc-core/pkg/api/store"
	"github.com/sprioc/sprioc-core/pkg/env"
	"github.com/sprioc/sprioc-core/pkg/model"
)

var dbname = env.Getenv("MONGODB_NAME", "morph")
var mongo = store.ConnectStore()

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
