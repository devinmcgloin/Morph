package store

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ds *MgoStore) GetCollection(IDs mgo.DBRef) (mgo.DBRef, error) {
	documents, err := get(ds, IDs)
	if err != nil {
		return mgo.DBRef{}, err
	}
	return documents.(mgo.DBRef), nil
}

func (ds *MgoStore) CreateCollection(Collection mgo.DBRef) error {
	return create(ds, "collections", Collection)
}

func (ds *MgoStore) DeleteCollection(ID mgo.DBRef) error {
	return delete(ds, ID)
}

func (ds *MgoStore) ModifyCollection(ID mgo.DBRef, diff bson.M) error {
	return modify(ds, ID, diff)
}

func (ds *MgoStore) AddUserToCollection(ID mgo.DBRef, ImageID mgo.DBRef) error {
	return modifyRef(ds, ID, ImageID, true)
}

func (ds *MgoStore) AddImageToCollection(ID mgo.DBRef, ImageID mgo.DBRef) error {
	return modifyRef(ds, ID, ImageID, true)
}

func (ds *MgoStore) DeleteImageFromCollection(ID mgo.DBRef, ImageID mgo.DBRef) error {
	return modifyRef(ds, ID, ImageID, false)
}

func (ds *MgoStore) DeleteUserFromCollection(ID mgo.DBRef, ImageID mgo.DBRef) error {
	return modifyRef(ds, ID, ImageID, false)
}

func (ds *MgoStore) FavoriteCollection(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "favorite", true)
}

func (ds *MgoStore) FollowCollection(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "follow", true)

}

func (ds *MgoStore) UnFavoriteCollection(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "favorite", false)

}

func (ds *MgoStore) UnFollowCollection(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "follow", false)
}
