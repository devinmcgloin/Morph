package sql

import (
	"log"

	"github.com/sprioc/composer/pkg/model"
)

func PatchImage(image model.Ref, changes map[string]interface{}) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	changes["id"] = image.Id

	_, err = tx.NamedExec(`UPDATE content.metadata (aperture, exposure_time, focal_length, iso, make, model, lens_model, lens_make, capture_time) 
	VALUES (:aperture, :exposure_time, :focal_length, :iso, :make, :model, :lens_model, :lens_make, :capture_time) WHERE image_id = :id`,
		changes)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func PatchUser(user model.Ref, changes map[string]interface{}) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	changes["id"] = user.Id

	_, err = tx.NamedExec(`UPDATE content.users () 
	VALUES () WHERE id = :id`,
		changes)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
