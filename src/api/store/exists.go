package store

import (
	"log"

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

func (ds *MgoStore) ExistsAlbumID(id bson.ObjectId) bool {
	return exists(ds, id, "albums")
}

func (ds *MgoStore) ExistsImageID(id bson.ObjectId) bool {
	return exists(ds, id, "images")
}

func (ds *MgoStore) ExistsUserID(id bson.ObjectId) bool {
	return exists(ds, id, "users")
}

func (ds *MgoStore) ExistsEventID(id bson.ObjectId) bool {
	return exists(ds, id, "events")
}

func (ds *MgoStore) ExistsCollectionID(id bson.ObjectId) bool {
	return exists(ds, id, "collections")
}

func exists(ds *MgoStore, id bson.ObjectId, collection string) bool {
	session := ds.session.Copy()
	defer session.Close()

	c := session.DB(dbname).C(collection)
	n, err := c.FindId(id).Count()
	if err != nil {
		log.Println(err)
		return false
	}
	if n > 0 {
		return true
	}
	return false
}
