package store

import (
	"errors"

	"github.com/sprioc/sprioc-core/pkg/model"
	"gopkg.in/mgo.v2/bson"
)

func (ds *MgoStore) GetCollection(colRef model.DBRef) (model.Collection, error) {
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

func (ds *MgoStore) CreateCollection(Collection model.DBRef) error {
	return create(ds, "collections", Collection)
}

func (ds *MgoStore) DeleteCollection(ID model.DBRef) error {
	return delete(ds, ID)
}

func (ds *MgoStore) ModifyCollection(ID model.DBRef, diff bson.M) error {
	return modify(ds, ID, diff)
}

func (ds *MgoStore) AddUserToCollection(ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, true)
}

func (ds *MgoStore) AddImageToCollection(ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, true)
}

func (ds *MgoStore) DeleteImageFromCollection(ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, false)
}

func (ds *MgoStore) DeleteUserFromCollection(ID model.DBRef, ImageID model.DBRef) error {
	return modifyRef(ds, ID, ImageID, false)
}

func (ds *MgoStore) FavoriteCollection(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", true)
}

func (ds *MgoStore) FollowCollection(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", true)

}

func (ds *MgoStore) UnFavoriteCollection(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", false)

}

func (ds *MgoStore) UnFollowCollection(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", false)
}
