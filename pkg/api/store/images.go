package store

import (
	"errors"

	"github.com/devinmcgloin/sprioc/pkg/model"
	"gopkg.in/mgo.v2/bson"
)

// TODO perhaps it would be good to do validation here.

func (ds *MgoStore) GetImage(imgRef model.DBRef) (model.Image, error) {
	session := ds.getSession()
	defer session.Close()

	var document model.Image

	c := session.DB(imgRef.Database).C(imgRef.Collection)

	err := c.Find(bson.M{"shortcode": imgRef.Shortcode}).One(&document)
	if err != nil {
		return model.Image{}, errors.New("Not found")
	}

	return document, nil
}

func (ds *MgoStore) CreateImage(image model.Image) error {
	return create(ds, "images", image)
}

func (ds *MgoStore) DeleteImage(ID model.DBRef) error {
	return delete(ds, ID)
}

func (ds *MgoStore) ModifyImage(ID model.DBRef, diff bson.M) error {
	return modify(ds, ID, diff)
}

func (ds *MgoStore) FavoriteImage(user model.DBRef, ID model.DBRef) error {
	return link(ds, user, ID, "favorite", true)
}

func (ds *MgoStore) FeatureImage(ID model.DBRef) error {
	err := modify(ds, ID, bson.M{"$set": bson.M{"featured": true}})
	if err != nil {
		return err
	}
	return nil
}

func (ds *MgoStore) UnFavoriteImage(user model.DBRef, ID model.DBRef) error {
	return link(ds, user, ID, "favorite", false)

}

func (ds *MgoStore) UnFeatureImage(ID model.DBRef) error {
	err := modify(ds, ID, bson.M{"$set": bson.M{"featured": false}})
	if err != nil {
		return err
	}
	return nil
}
