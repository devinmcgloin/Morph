package store

import (
	"errors"

	"github.com/sprioc/sprioc-core/pkg/model"
	"gopkg.in/mgo.v2/bson"
)

func (ds *MgoStore) GetAlbum(albRef model.DBRef) (model.Album, error) {
	session := ds.getSession()
	defer session.Close()

	var document model.Album

	c := session.DB(albRef.Database).C(albRef.Collection)

	err := c.Find(bson.M{"shortcode": albRef.Shortcode}).One(&document)
	if err != nil {
		return model.Album{}, errors.New("Not found")
	}

	return document, nil
}

func (ds *MgoStore) CreateAlbum(Album model.DBRef) error {
	return create(ds, "Albums", Album)
}

func (ds *MgoStore) DeleteAlbum(ID model.DBRef) error {
	return delete(ds, ID)
}

func (ds *MgoStore) ModifyAlbum(ID model.DBRef, diff bson.M) error {
	return modify(ds, ID, diff)
}

func (ds *MgoStore) AddImageToAlbum(ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, true)
}

func (ds *MgoStore) DeleteImageFromAlbum(ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, false)
}

func (ds *MgoStore) FavoriteAlbum(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", true)
}

func (ds *MgoStore) FollowAlbum(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", true)
}

func (ds *MgoStore) UnFavoriteAlbum(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", false)
}

func (ds *MgoStore) UnFollowAlbum(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", false)
}
