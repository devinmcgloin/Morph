package ref

import (
	"github.com/devinmcgloin/sprioc/src/api/store"
	"github.com/devinmcgloin/sprioc/src/env"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dbname = env.Getenv("MONGODB_NAME", "morph")
var mongo = store.ConnectStore()

func GetImageRef(ID string) mgo.DBRef {
	return mgo.DBRef{Database: dbname, Collection: "images", Id: bson.ObjectId(ID)}
}
func GetCollectionRef(ID string) mgo.DBRef {
	return mgo.DBRef{Database: dbname, Collection: "collections", Id: bson.ObjectId(ID)}
}
func GetAlbumRef(ID string) mgo.DBRef {
	return mgo.DBRef{Database: dbname, Collection: "albums", Id: bson.ObjectId(ID)}
}
func GetUserRef(ID bson.ObjectId) mgo.DBRef {
	return mgo.DBRef{Database: dbname, Collection: "users", Id: ID}
}
