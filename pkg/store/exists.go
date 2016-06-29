package store

import (
	"log"

	"github.com/sprioc/sprioc-core/pkg/model"

	"gopkg.in/mgo.v2/bson"
)

func ExistsEmail(email string) bool {
	session := mongo.session.Copy()
	defer session.Close()

	c := session.DB(dbname).C("users")
	n, err := c.Find(bson.M{"email": email}).Count()
	if err != nil {
		log.Println(err)
		return false
	}
	if n > 0 {
		return true
	}
	return false
}

func ExistsAlbumID(id string) bool {
	return exists(model.DBRef{Database: dbname, Collection: "albums"}, bson.M{"shortcode": id})
}

func ExistsImageID(id string) bool {
	return exists(model.DBRef{Database: dbname, Collection: "images"}, bson.M{"shortcode": id})
}

func ExistsUserID(id string) bool {
	return exists(model.DBRef{Database: dbname, Collection: "users"}, bson.M{"shortcode": id})
}

func ExistsEventID(id string) bool {
	return exists(model.DBRef{Database: dbname, Collection: "events"}, bson.M{"shortcode": id})
}

func ExistsCollectionID(id string) bool {
	return exists(model.DBRef{Database: dbname, Collection: "collections"}, bson.M{"shortcode": id})
}

func exists(ref model.DBRef, query bson.M) bool {
	session := mongo.session.Copy()
	defer session.Close()

	c := session.DB(ref.Database).C(ref.Collection)
	n, err := c.Find(query).Count()
	if err != nil {
		log.Println(err)
		return false
	}
	if n > 0 {
		return true
	}
	return false
}
