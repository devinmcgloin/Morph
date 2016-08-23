package store

import (
	"log"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/qmgo"

	"gopkg.in/mgo.v2/bson"
)

func init() {
	ValidateConnection()
	mongo = ConnectStore()
}

var mongo *MgoStore

// TODO need to check if modification already exists and that types are correct.
// Bools should be bools. Only need to worry about multiple requests when
// working with lists.

// TODO should say something if the operation does not do anything.

func Get(ID model.DBRef, container interface{}) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(ID.Database).C(ID.Collection)

	err := c.Find(bson.M{"shortcode": ID.Shortcode}).One(container)
	if err != nil {
		log.Println(ID, err)
		return err
	}

	return nil
}

func GetAll(collection string, filter bson.M, dest interface{}) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C(collection)

	err := c.Find(filter).All(dest)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetImages(filter qmgo.ImageSearch, dest interface{}) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C("images")

	err := c.Find(filter).All(dest)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func Create(collection string, document interface{}) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C(collection)

	err := c.Insert(document)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Delete(ID model.DBRef) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C(ID.Collection)

	err := c.Remove(bson.M{"shortcode": ID.Shortcode})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Modify(ID model.DBRef, changes bson.M) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C(ID.Collection)

	err := c.Update(bson.M{"shortcode": ID.Shortcode}, changes)
	if err != nil {
		log.Println(err, ID, changes)
		return err
	}
	return nil
}
