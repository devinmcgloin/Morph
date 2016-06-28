package refs

import (
	"github.com/sprioc/sprioc-core/pkg/env"
	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/store"
)

var dbname = env.Getenv("MONGODB_NAME", "sprioc")
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
