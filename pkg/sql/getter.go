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
	images := []string{}
	err = db.Select(&images, `
	SELECT shortcode
	FROM content.images AS images
	WHERE owner_id = $1`, user.Id)
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}
	user.Images = images
	return user, nil
}

func GetImage(i string) (model.Image, error) {
	img := model.Image{}
	err := db.Get(&img, `
	SELECT images.id, shortcode, publish_time, images.last_modified,
		owner_id, users.username, images.featured, images.downloads, images.views,
		aperture, exposure_time, focal_length, iso, make, model, lens_make, lens_model,
		pixel_xd, pixel_yd, capture_time
	FROM content.images AS images
		JOIN content.users AS users ON owner_id = users.id
		JOIN content.image_metadata AS metadata ON image_id = images.id
	WHERE shortcode = $1`, i)
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}
	return img, nil
}
func GetRecentImages(limit int) ([]model.Image, error) {
	imgs := []model.Image{}
	err := db.Select(&imgs,
		"SELECT * FROM content.images ORDER BY publish_time DESC LIMIT $1",
		limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return imgs, nil
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
