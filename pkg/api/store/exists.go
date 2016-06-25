package store

import (
	"log"

	"github.com/devinmcgloin/sprioc/pkg/model"

	"gopkg.in/mgo.v2/bson"
)

func (ds *MgoStore) ExistsUserName(userName string) bool {
	session := ds.session.Copy()
	defer session.Close()

	c := session.DB(dbname).C("users")
	n, err := c.Find(bson.M{"username": userName}).Count()
	if err != nil {
		log.Println(err)
		return false
	}
	if n > 0 {
		return true
	}
	return false
}

func (ds *MgoStore) ExistsEmail(email string) bool {
	session := ds.session.Copy()
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

func (ds *MgoStore) ExistsAlbumID(id string) bool {
	return exists(ds, model.DBRef{Database: dbname, Collection: "albums"}, bson.M{"shortcode": id})
}

func (ds *MgoStore) ExistsImageID(id string) bool {
	return exists(ds, model.DBRef{Database: dbname, Collection: "images"}, bson.M{"shortcode": id})
}

func (ds *MgoStore) ExistsUserID(id string) bool {
	return exists(ds, model.DBRef{Database: dbname, Collection: "users"}, bson.M{"shortcode": id})
}

func (ds *MgoStore) ExistsEventID(id string) bool {
	return exists(ds, model.DBRef{Database: dbname, Collection: "events"}, bson.M{"shortcode": id})
}

func (ds *MgoStore) ExistsCollectionID(id string) bool {
	return exists(ds, model.DBRef{Database: dbname, Collection: "collections"}, bson.M{"shortcode": id})
}

func exists(ds *MgoStore, ref model.DBRef, query bson.M) bool {
	session := ds.session.Copy()
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
