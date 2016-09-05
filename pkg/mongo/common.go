package mongo

import (
	"fmt"
	"log"
	"os"

	"github.com/sprioc/composer/pkg/model"

	"gopkg.in/mgo.v2/bson"
)

var (
	database string = os.Getenv("MONGODB_NAME")
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

func Get(ID model.Ref, container interface{}) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(database).C(string(ID.Collection))

	err := c.Find(bson.M{"shortcode": ID.ShortCode}).One(container)
	if err != nil {
		log.Printf("%+v, %v", ID, err)
		return err
	}

	return nil
}

func GetAll(collection string, filter bson.M, dest interface{}) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C(collection)
	var err error
	switch t := dest.(type) {
	case []model.Image:
		err = c.Find(filter).All(dest.([]model.Image))
	case []model.Collection:
		err = c.Find(filter).All(dest.([]model.Collection))
	case []model.User:
		err = c.Find(filter).All(dest.([]model.User))
	default:
		return fmt.Errorf("Invalid dest type %s", t)
	}

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// func SearchImages(filter qmgo.ImageSearch, dest interface{}) error {
// 	session := mongo.getSession()
// 	defer session.Close()
//
// 	c := session.DB(dbname).C("images")
//
// 	var err error
// 	if len(filter.Sort) == 0 {
// 		err = c.Find(filter).All(dest)
// 	} else {
// 		err = c.Find(filter).Sort(filter.Sort...).All(dest)
// 	}
//
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
//
// 	return nil
// }

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

func Delete(ID model.Ref) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C(string(ID.Collection))

	err := c.Remove(bson.M{"shortcode": ID.ShortCode})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Modify(ID model.Ref, changes bson.M) error {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(dbname).C(string(ID.Collection))

	err := c.Update(bson.M{"shortcode": ID.ShortCode}, changes)
	if err != nil {
		log.Println(err, ID, changes)
		return err
	}
	return nil
}
