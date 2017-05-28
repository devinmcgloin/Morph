package sql

import (
	"log"

	"github.com/lib/pq"
	"github.com/sprioc/composer/pkg/model"
)

// GetUser returns the fields of a user row into a User struct, including image references.
func GetUser(u int64) (model.User, error) {
	user := model.User{}
	err := db.Get(&user, "SELECT * FROM content.users WHERE id = $1", u)
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

	favorites := []string{}
	err = db.Select(&favorites, `
	SELECT shortcode
	FROM content.user_favorites
	JOIN content.images ON content.user_favorites.image_id = content.images.id
	WHERE user_id = $1`, user.Id)
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}
	user.Favorites = favorites

	return user, nil
}

// GetImage takes an image ID and returns a image row into a Image struct including metadata
// and user data.
func GetImage(i int64) (model.Image, error) {
	img := model.Image{}
	err := db.Get(&img, `
	SELECT images.id, shortcode, publish_time, images.last_modified,
		owner_id, users.username, images.featured, images.downloads, images.views,
		aperture, exposure_time, focal_length, iso, make, model, lens_make, lens_model,
		pixel_xd, pixel_yd, capture_time
	FROM content.images AS images
		JOIN content.users AS users ON owner_id = users.id
		JOIN content.image_metadata AS metadata ON image_id = images.id
	WHERE images.id = $1`, i)
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}

	var tags []string
	err = db.Select(&tags, `
	SELECT description FROM content.image_tags
	JOIN content.image_tag_bridge ON content.image_tags.id = content.image_tag_bridge.tag_id
	WHERE image_tag_bridge.image_id = $1;
	`, img.Id)
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}
	img.Tags = tags
	return img, nil
}

func GetImageID(i string) (int64, error) {
	var iID int64

	err := db.Get(&iID, "SELECT id FROM content.images WHERE shortcode = $1", i)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return iID, nil
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

func GetUserFollowed(userID int64) ([]model.User, error) {
	users := []model.User{}
	err := db.Select(&users,
		`
	SELECT id, username, email, name, bio, url, password, salt, featured, admin, 
		views, created_at, last_modified
	FROM content.user_follows AS follows 
		JOIN content.users AS users ON id = follows.followed_id
	WHERE follows.user_id = $1
		`,
		userID)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.User{}, err
	}
	return users, nil
}

func GetUserFavorites(userID int64) ([]model.Image, error) {
	imgs := []model.Image{}
	err := db.Select(&imgs,
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
		userID)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return imgs, nil
}

func GetUserImages(userID int64) ([]model.Image, error) {
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
	WHERE owner_id = $1
	ORDER BY publish_time DESC
		`,
		userID)
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

func GetTagRef(t string) (model.Ref, error) {
	ref := model.Ref{Collection: model.Tags}
	err := db.Get(&ref.Id, "SELECT id FROM content.image_tags WHERE description = $1", t)
	if err != nil {
		log.Println(err)
		rows := db.QueryRow("INSERT INTO content.image_tags (description) VALUES ($1) RETURNING id", t)
		err = rows.Scan(&ref.Id)
		if err != nil {
			return model.Ref{}, err
		}
		return ref, nil
	}
	return ref, nil
}

func GetImageRef(i string) (model.Ref, error) {
	ref := model.Ref{Collection: model.Images}
	err := db.Get(&ref.Id, "SELECT id FROM content.images WHERE shortcode = $1", i)
	if err != nil {
		log.Println(err)
		return model.Ref{}, err
	}
	return ref, nil
}

func GetUserRef(u string) (model.Ref, error) {
	ref := model.Ref{Collection: model.Users}
	err := db.Get(&ref.Id, "SELECT id FROM content.users WHERE username = $1", u)
	if err != nil {
		log.Println(err)
		return model.Ref{}, err
	}
	return ref, nil
}
