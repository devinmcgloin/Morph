package store

import (
	"errors"

	"github.com/devinmcgloin/sprioc/src/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ds *MgoStore) GetAlbum(albRef mgo.DBRef) (model.Album, error) {
	session := ds.getSession()
	defer session.Close()

	var document model.Album

	err := session.FindRef(&albRef).One(&document)
	if err != nil {
		return model.Album{}, errors.New("Not found")
	}

	return document, nil
}

func (ds *MgoStore) CreateAlbum(Album mgo.DBRef) error {
	return create(ds, "Albums", Album)
}

func (ds *MgoStore) DeleteAlbum(ID mgo.DBRef) error {
	return delete(ds, ID)
}

func (ds *MgoStore) ModifyAlbum(ID mgo.DBRef, diff bson.M) error {
	return modify(ds, ID, diff)
}

func (ds *MgoStore) AddImageToAlbum(ID mgo.DBRef, ImageID mgo.DBRef) error {
	return modifyRef(ds, ID, ImageID, true)
}

func (ds *MgoStore) DeleteImageFromAlbum(ID mgo.DBRef, ImageID mgo.DBRef) error {
	return modifyRef(ds, ID, ImageID, false)
}

func (ds *MgoStore) FavoriteAlbum(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "favorite", true)
}

func (ds *MgoStore) FollowAlbum(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "follow", true)
}

func (ds *MgoStore) UnFavoriteAlbum(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "favorite", false)
}

func (ds *MgoStore) UnFollowAlbum(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "follow", false)
}
