package store

import (
	"errors"
	"log"

	"github.com/devinmcgloin/sprioc/pkg/model"
	"gopkg.in/mgo.v2/bson"
)

func (ds *MgoStore) GetUser(userRef model.DBRef) (model.User, error) {
	session := ds.getSession()
	defer session.Close()

	var document model.User

	c := session.DB(userRef.Database).C(userRef.Collection)

	err := c.Find(bson.M{"shortcode": userRef.Shortcode}).One(&document)
	if err != nil {
		return model.User{}, errors.New("Not found")
	}

	return document, nil
}

func (ds *MgoStore) CreateUser(user model.User) error {
	return create(ds, "users", user)
}

func (ds *MgoStore) DeleteUser(userRef model.DBRef) error {
	return delete(ds, userRef)
}

func (ds *MgoStore) ModifyUser(userRef model.DBRef, diff bson.M) error {
	return modify(ds, userRef, diff)
}

func (ds *MgoStore) ModifyAvatar(userRef model.DBRef, urls model.ImgSource) error {
	return modify(ds, userRef, bson.M{"$set": bson.M{"avatar_url": urls}})
}

func (ds *MgoStore) FavoriteUser(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", true)
}

func (ds *MgoStore) FollowUser(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", true)

}

func (ds *MgoStore) UnFavoriteUser(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", false)
}

func (ds *MgoStore) UnFollowUser(actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", false)
}

func (ds *MgoStore) GetByUserName(username string) (model.User, error) {
	session := ds.getSession()
	defer session.Close()

	var document model.User

	c := session.DB(dbname).C("users")

	err := c.Find(bson.M{"shortcode": username}).One(&document)
	if err != nil {
		log.Println(err)
		return model.User{}, errors.New("Not Found")
	}

	return document, nil

}
