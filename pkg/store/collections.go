package store

import (
	"errors"

	"github.com/sprioc/sprioc-core/pkg/model"
	"gopkg.in/mgo.v2/bson"
)

func GetCollection(ds *MgoStore, colRef model.DBRef) (model.Collection, error) {
	session := ds.getSession()
	defer session.Close()

	var document model.Collection

	c := session.DB(colRef.Database).C(colRef.Collection)

	err := c.Find(bson.M{"shortcode": colRef.Shortcode}).One(&document)
	if err != nil {
		return model.Collection{}, errors.New("Not found")
	}

	return document, nil
}

func CreateCollection(ds *MgoStore, Collection model.DBRef) error {
	return create(ds, "collections", Collection)
}

func DeleteCollection(ds *MgoStore, ID model.DBRef) error {
	return delete(ds, ID)
}

func ModifyCollection(ds *MgoStore, ID model.DBRef, diff bson.M) error {
	return modify(ds, ID, diff)
}

func AddUserToCollection(ds *MgoStore, ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, true)
}

func AddImageToCollection(ds *MgoStore, ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, true)
}

func DeleteImageFromCollection(ds *MgoStore, ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, false)
}

func DeleteUserFromCollection(ds *MgoStore, ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, false)
}

func FavoriteCollection(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", true)
}

func FollowCollection(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", true)

}

func UnFavoriteCollection(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", false)

}

func UnFollowCollection(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", false)
}
