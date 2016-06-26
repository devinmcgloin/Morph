package store

import (
	"errors"

	"github.com/sprioc/sprioc-core/pkg/model"
	"gopkg.in/mgo.v2/bson"
)

func GetUser(ds *MgoStore, userRef model.DBRef) (model.User, error) {
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

func CreateUser(ds *MgoStore, user model.User) error {
	return create(ds, "users", user)
}

func DeleteUser(ds *MgoStore, userRef model.DBRef) error {
	return delete(ds, userRef)
}

func ModifyUser(ds *MgoStore, userRef model.DBRef, diff bson.M) error {
	return modify(ds, userRef, diff)
}

func ModifyAvatar(ds *MgoStore, userRef model.DBRef, avatar model.ImgSource) error {
	return modify(ds, userRef, bson.M{"$set": bson.M{"avatar_url": avatar}})
}

func FavoriteUser(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", true)
}

func FollowUser(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", true)

}

func UnFavoriteUser(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "favorite", false)
}

func UnFollowUser(ds *MgoStore, actor model.DBRef, recipient model.DBRef) error {
	return link(ds, actor, recipient, "follow", false)
}
