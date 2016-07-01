package store

import (
	"log"

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
	return Exists("albums", bson.M{"shortcode": id})
}

func ExistsImageID(id string) bool {
	return Exists("images", bson.M{"shortcode": id})
}

func ExistsUserID(id string) bool {
	return Exists("users", bson.M{"shortcode": id})
}

func ExistsEventID(id string) bool {
	return Exists("events", bson.M{"shortcode": id})
}

func ExistsCollectionID(id string) bool {
	return Exists("collections", bson.M{"shortcode": id})
}

func Exists(collection string, query bson.M) bool {
	session := mongo.session.Copy()
	defer session.Close()

	c := session.DB(dbname).C(collection)
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
