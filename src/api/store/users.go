package store

import (
	"errors"

	"github.com/devinmcgloin/sprioc/src/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ds *MgoStore) GetUser(userRef mgo.DBRef) (model.User, error) {
	documents, err := get(ds, userRef)
	if err != nil {
		return model.User{}, err
	}
	return documents.(model.User), nil
}

func (ds *MgoStore) CreateUser(user model.User) error {
	return create(ds, "users", user)
}

func (ds *MgoStore) DeleteUser(userRef mgo.DBRef) error {
	return delete(ds, userRef)
}

func (ds *MgoStore) ModifyUser(userRef mgo.DBRef, diff bson.M) error {
	return modify(ds, userRef, diff)
}

func (ds *MgoStore) ModifyAvatar(userRef mgo.DBRef, url model.URL) error {
	return modify(ds, userRef, bson.M{"$set": bson.M{"avatar_url": url}})
}

func (ds *MgoStore) FavoriteUser(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "favorite", true)
}

func (ds *MgoStore) FollowUser(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "follow", true)

}

func (ds *MgoStore) UnFavoriteUser(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "favorite", false)
}

func (ds *MgoStore) UnFollowUser(actor mgo.DBRef, recipient mgo.DBRef) error {
	return link(ds, actor, recipient, "follow", false)
}

func (ds *MgoStore) GetByUserName(username model.UserName) (model.User, error) {
	session := ds.getSession()
	defer session.Close()

	var document model.User

	c := session.DB(dbname).C("users")

	err := c.Find(bson.M{"username": username}).One(&document)
	if err != nil {
		return model.User{}, errors.New("Not Found")
	}

	return document, nil

}
