package sql

import (
	"log"

	"github.com/sprioc/composer/pkg/model"
)

// CreateImage stores the image data in the database under the given user.
// Currently does not set the metadata or db interal state.
func CreateImage(image model.Image) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}
	var id int64
	rows, err := tx.NamedQuery(`
	INSERT INTO content.images(owner_id, shortcode)
	VALUES(:owner_id, :shortcode) RETURNING id;`,
		image)
	if err != nil {
		log.Println(err)
		return err
	}

	for rows.Next() {
		rows.Scan(&id)
	}

	image.Id = id

	_, err = tx.NamedExec(`
	INSERT INTO content.image_metadata(image_id, aperture, exposure_time, 
	focal_length, iso, make, model, lens_make, lens_model, pixel_xd, 
	pixel_yd, capture_time) VALUES (:id, :aperture, :exposure_time,
	:focal_length, :iso, :make, :model, :lens_make, :lens_model, :pixel_xd, 
	:pixel_yd, :capture_time);
	`, image)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func CreateUser(user model.User) error {
	log.Println(user)
	_, err := db.NamedExec(`
	INSERT INTO content.users(username, email, password, salt)
	VALUES(:username, :email, :password, :salt);`,
		user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
