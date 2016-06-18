package store

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

func (ds *MgoStore) ExistsUserName(userName string) bool {
	session := ds.session.Copy()
	defer session.Close()

	c := session.DB("morph").C("users")
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

func (ds *MgoStore) ExistsAlbumShortCode(shortCode string) bool {
	session := ds.session.Copy()
	defer session.Close()

	c := session.DB("morph").C("album")
	n, err := c.Find(bson.M{"shortcode": shortCode}).Count()
	if err != nil {
		log.Println(err)
		return false
	}
	if n > 0 {
		return true
	}
	return false
}

func (ds *MgoStore) ExistsImageShortCode(shortCode string) bool {
	session := ds.session.Copy()
	defer session.Close()

	c := session.DB("morph").C("images")
	n, err := c.Find(bson.M{"shortcode": shortCode}).Count()
	if err != nil {
		log.Println(err)
		return false
	}
	if n > 0 {
		return true
	}
	return false
}

func (ds *MgoStore) ExistsUser(provider string, providerID string) bool {
	session := ds.session.Copy()
	defer session.Close()

	c := session.DB("morph").C("users")
	n, err := c.Find(bson.M{"provider": provider, "provider_id": providerID}).Count()
	if err != nil {
		log.Println(err)
		return false
	}
	if n > 0 {
		return true
	}
	return false
}
