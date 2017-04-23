package sql

import (
	"log"

	"github.com/sprioc/composer/pkg/model"
)

func GetUser(u string) (model.User, error) {
	user := model.User{}
	err := db.Get(&user, "SELECT * FROM content.users WHERE username = $1", u)
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}
	return user, nil
}

func GetImage(i string) (model.Image, error) {
	img := model.Image{}
	err := db.Get(&img, "SELECT * FROM content.images WHERE shortcode = $1", i)
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}
	return img, nil
}
