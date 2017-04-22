package sql

import (
	"log"

	"github.com/sprioc/composer/pkg/model"
)

func GetUser(u model.UserReference) (model.User, error) {
	user := model.User{}
	err := db.Get(&user, "SELECT * FROM content.users WHERE id = ?", u)
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}
	return user, nil
}

func GetImage(i model.ImageReference) (model.Image, error) {
	img := model.Image{}
	err := db.Get(&img, "SELECT * FROM content.images WHERE id = ?", i)
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}
	return img, nil
}
