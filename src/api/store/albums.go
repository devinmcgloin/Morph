package store

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ds *MgoStore) GetAlbum(IDs mgo.DBRef) (mgo.DBRef, error) {
	documents, err := get(ds, IDs)
	if err != nil {
		return mgo.DBRef{}, err
	}
	return documents.(mgo.DBRef), nil
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
