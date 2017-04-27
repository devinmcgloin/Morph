package sql

import (
	"log"

	"github.com/lib/pq"
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

func GetFeaturedImages(limit int) ([]model.Image, error) {
	imgs := []model.Image{}
	err := db.Select(&imgs,
		"SELECT * FROM content.images WHERE featured = TRUE ORDER BY publish_time DESC LIMIT $1",
		limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return imgs, nil
}
