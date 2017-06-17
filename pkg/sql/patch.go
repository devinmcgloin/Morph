package sql

import (
	"fmt"
	"log"

	"github.com/sprioc/composer/pkg/model"
)

func PatchImage(image model.Ref, changes map[string]interface{}) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.Exec("DELETE FROM content.image_tag_bridge WHERE image_id = $1", image.Id)
	if err != nil {
		log.Println(err)
		return err
	}

	for key, val := range changes {
		if key == "tags" {
			tags := val.([]string)
			var tagID int
			for _, tag := range tags {
				tx.Get(&tagID, "SELECT (id) WHERE description = $1",
					tag)
				_, err = tx.Exec(`INSERT INTO content.image_tags_bridge (image_id, tag_id)
			VALUES ($1, $2)`, image.Id, tagID)
				if err != nil {
					log.Println(err)
					return err
				}
			}
		} else {
			_, err = tx.Exec(`UPDATE content.metadata ($1) VALUES ($2) WHERE image_id = $3`, val, key, image.Id)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	err = tx.Commit()
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

	for key, val := range changes {
		log.Printf("UPDATE content.users set %s = %s WHERE id = %d", key, val, user.Id)
		stmt := fmt.Sprintf("UPDATE content.users SET %s = $1 WHERE id = $2", key)
		_, err = tx.Exec(stmt, val, user.Id)
		if err != nil {
			log.Println(err)
			return err
		}

	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
