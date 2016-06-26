package store

import (
	"errors"

	"github.com/sprioc/sprioc-core/pkg/model"
	"gopkg.in/mgo.v2/bson"
)

func GetAlbum(ds *MgoStore, albRef model.DBRef) (model.Album, error) {
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

func CreateAlbum(ds *MgoStore, Album model.DBRef) error {
	return create(ds, "Albums", Album)
}

func DeleteAlbum(ds *MgoStore, ID model.DBRef) error {
	return delete(ds, ID)
}

func ModifyAlbum(ds *MgoStore, ID model.DBRef, diff bson.M) error {
	return modify(ds, ID, diff)
}

func AddImageToAlbum(ds *MgoStore, ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, true)
}

func DeleteImageFromAlbum(ds *MgoStore, ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, false)
}

func FavoriteAlbum(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", true)
}

func FollowAlbum(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", true)
}

func UnFavoriteAlbum(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", false)
}

func UnFollowAlbum(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", false)
}
