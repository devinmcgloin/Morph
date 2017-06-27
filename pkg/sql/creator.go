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

	rows.Close()

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

	_, err = tx.Exec(`
	INSERT INTO content.image_geo (image_id, loc, dir) 
	VALUES ($1, ST_GeographyFromText('SRID=4326;POINT($2 $3)'), $4);
	`, image.Id, image.Location.Coordinates[0], image.Location.Coordinates[1],
		image.ImgDirection)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_edit(user_id, o_id, type) VALUES (:owner_id, :id, 'image');
	`, image)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_delete(user_id, o_id, type) VALUES (:owner_id, :id, 'image');
	`, image)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_view(user_id, o_id, type) VALUES (-1, :id, 'image');
	`, image)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			return err
		}
		return err
	}
	return nil
}

func CreateUser(user model.User) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.NamedExec(`
	INSERT INTO content.users(username, email, password, salt)
	VALUES(:username, :email, :password, :salt);`,
		user)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(user)
	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_edit(user_id, o_id, type) VALUES (:id, :id, 'user')`, user)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_delete(user_id, o_id, type) VALUES (:id, :id, 'user')`, user)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.NamedQuery(`
	INSERT INTO permissions.can_view(user_id, o_id, type) VALUES (-1, :id, 'user')`, user)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
