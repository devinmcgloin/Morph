package store

import (
	"errors"
	"log"

	"github.com/sprioc/sprioc-core/pkg/model"
	"gopkg.in/mgo.v2/bson"
)

// TODO perhaps it would be good to do validation here.

func GetImage(ds *MgoStore, imgRef model.DBRef) (model.Image, error) {
	session := ds.getSession()
	defer session.Close()

	log.Println(imgRef)

	var document model.Image

	c := session.DB(imgRef.Database).C(imgRef.Collection)

	err := c.Find(bson.M{"shortcode": imgRef.Shortcode}).One(&document)
	if err != nil {
		log.Println(err)
		return model.Image{}, errors.New("Not found")
	}
	return document, nil
}

func CreateImage(ds *MgoStore, image model.Image) error {
	return create(ds, "images", image)
}

func DeleteImage(ds *MgoStore, ID model.DBRef) error {
	return delete(ds, ID)
}

func ModifyImage(ds *MgoStore, ID model.DBRef, diff bson.M) error {
	return modify(ds, ID, diff)
}

func FavoriteImage(ds *MgoStore, user model.DBRef, ID model.DBRef) error {
	return link(ds, user, ID, "favorite", true)
}

func FeatureImage(ds *MgoStore, ID model.DBRef) error {
	err := modify(ds, ID, bson.M{"$set": bson.M{"featured": true}})
	if err != nil {
		return err
	}
	return nil
}

func UnFavoriteImage(ds *MgoStore, user model.DBRef, ID model.DBRef) error {
	return link(ds, user, ID, "favorite", false)

}

func UnFeatureImage(ds *MgoStore, ID model.DBRef) error {
	err := modify(ds, ID, bson.M{"$set": bson.M{"featured": false}})
	if err != nil {
		return err
	}
	return nil
}
