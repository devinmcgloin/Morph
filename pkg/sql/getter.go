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
		`
	SELECT images.id, shortcode, publish_time, images.last_modified,
		owner_id, users.username, images.featured, images.downloads, images.views,
		aperture, exposure_time, focal_length, iso, make, model, lens_make, lens_model,
		pixel_xd, pixel_yd, capture_time
	FROM content.images AS images
		JOIN content.users AS users ON owner_id = users.id
		JOIN content.image_metadata AS metadata ON image_id = images.id
	ORDER BY publish_time DESC LIMIT $1
		`,
		limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return imgs, nil
}

func GetUserFavorites(username string) ([]model.Image, error) {
	imgs := []model.Image{}
	var owner_id int64
	err := db.Get(&owner_id, "SELECT id FROM content.users WHERE username = $1", username)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	err = db.Select(&imgs,
		`
	SELECT images.id, shortcode, publish_time, images.last_modified,
		owner_id, users.username, images.featured, images.downloads, images.views,
		aperture, exposure_time, focal_length, iso, make, model, lens_make, lens_model,
		pixel_xd, pixel_yd, capture_time
	FROM content.user_favorites AS favs
		JOIN content.images AS images ON favs.image_id = images.id 
		JOIN content.users AS users ON owner_id = users.id
		JOIN content.image_metadata AS metadata ON metadata.image_id = images.id
	WHERE favs.user_id = $1
	ORDER BY publish_time DESC
		`,
		owner_id)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return imgs, nil
}

func GetUserImages(username string) ([]model.Image, error) {
	imgs := []model.Image{}
	var owner_id int64
	err := db.Get(&owner_id, "SELECT id FROM content.users WHERE username = $1", username)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	err = db.Select(&imgs,
		`
	SELECT images.id, shortcode, publish_time, images.last_modified,
		owner_id, users.username, images.featured, images.downloads, images.views,
		aperture, exposure_time, focal_length, iso, make, model, lens_make, lens_model,
		pixel_xd, pixel_yd, capture_time
	FROM content.images AS images
		JOIN content.users AS users ON owner_id = users.id
		JOIN content.image_metadata AS metadata ON image_id = images.id
	WHERE owner_id = $1
	ORDER BY publish_time DESC
		`,
		owner_id)
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
		`
	SELECT images.id, shortcode, publish_time, images.last_modified,
		owner_id, users.username, images.featured, images.downloads, images.views,
		aperture, exposure_time, focal_length, iso, make, model, lens_make, lens_model,
		pixel_xd, pixel_yd, capture_time
	FROM content.images AS images
		JOIN content.users AS users ON owner_id = users.id
		JOIN content.image_metadata AS metadata ON image_id = images.id
	WHERE featured = TRUE 
	ORDER BY publish_time DESC LIMIT $1
		`,
		limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return imgs, nil
}
