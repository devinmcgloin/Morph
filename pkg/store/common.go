package store

import (
	"errors"
	"log"
	"strings"

	"github.com/sprioc/sprioc-core/pkg/model"

	"gopkg.in/mgo.v2/bson"
)

func init() {
	relationTerms = make(map[string]string)
	relationTerms["follow_forward"] = "followes"
	relationTerms["follow_backward"] = "followed_by"

	relationTerms["favorite_forward"] = "favorites"
	relationTerms["favorite_backward"] = "favorited_by"

	relationTerms["collection_forward"] = "images"
	relationTerms["collection_backward"] = "collections"
}

var relationTerms map[string]string

var mongo = ConnectStore()

// TODO need to check if modification already exists and that types are correct.
// Bools should be bools. Only need to worry about multiple requests when
// working with lists.

// TODO should say something if the operation does not do anything.

func Get(ID model.DBRef) (interface{}, error) {
	session := mongo.getSession()
	defer session.Close()

	var document interface{}

	c := session.DB(ID.Database).C(ID.Collection)

	err := c.Find(bson.M{"shortcode": ID.Shortcode}).One(&document)
	if err != nil {
		return nil, errors.New("Not found")
	}

	return document, nil
}

func Create(collection string, document interface{}) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C(collection)

	err := c.Insert(document)
	if err != nil {
		log.Println(err)
		return errors.New("Unable to add document to DB")
	}
	return nil
}

func Delete(ID model.DBRef) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C(ID.Collection)

	return c.Remove(bson.M{"shortcode": ID.Shortcode})
}

func Modify(ID model.DBRef, changes bson.M) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C(ID.Collection)

	return c.Update(bson.M{"shortcode": ID.Shortcode}, changes)
}

func Link(actor model.DBRef, recipient model.DBRef, relation string) error {
	return errors.New("Not Implemented")
}

func Unlink(actor model.DBRef, recipient model.DBRef, relation string) error {
	return errors.New("Not Implemented")
}

func in(contentType string, opts []string) bool {
	for _, opt := range opts {
		if strings.Compare(contentType, opt) == 0 {
			return true
		}
	}
	return false
}
