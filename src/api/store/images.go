package store

import (
	"errors"

	"github.com/devinmcgloin/sprioc/src/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// TODO perhaps it would be good to do validation here.

func (ds *MgoStore) GetImage(imgRef mgo.DBRef) (model.Image, error) {
	session := ds.getSession()
	defer session.Close()

	var document model.Image

	err := session.FindRef(&imgRef).One(&document)
	if err != nil {
		return model.Image{}, errors.New("Not found")
	}

	return document, nil
}

func (ds *MgoStore) CreateImage(image model.Image) error {
	return create(ds, "images", image)
}

func (ds *MgoStore) DeleteImage(ID mgo.DBRef) error {
	return delete(ds, ID)
}

func (ds *MgoStore) ModifyImage(ID mgo.DBRef, diff bson.M) error {
	return modify(ds, ID, diff)
}

func (ds *MgoStore) FavoriteImage(user mgo.DBRef, ID mgo.DBRef) error {
	return link(ds, user, ID, "favorite", true)
}

func (ds *MgoStore) FeatureImage(ID mgo.DBRef) error {
	err := modify(ds, ID, bson.M{"$set": bson.M{"featured": true}})
	if err != nil {
		return err
	}
	return nil
}

func (ds *MgoStore) UnFavoriteImage(user mgo.DBRef, ID mgo.DBRef) error {
	return link(ds, user, ID, "favorite", false)

}

func (ds *MgoStore) UnFeatureImage(ID mgo.DBRef) error {
	err := modify(ds, ID, bson.M{"$set": bson.M{"featured": false}})
	if err != nil {
		return err
	}
	return nil
}
