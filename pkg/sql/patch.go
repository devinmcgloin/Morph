package sql

import (
	"fmt"
	"log"

	"github.com/devinmcgloin/fokal/pkg/model"
)

func PatchImage(image model.Ref, changes map[string]interface{}) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	for key, val := range changes {
		if key == "tags" {
			_, err = tx.Exec("DELETE FROM content.image_tag_bridge WHERE image_id = $1;", image.Id)
			if err != nil {
				log.Println(err)
				return err
			}

			tags := val.([]interface{})
			var tagID int
			for _, tag := range tags {
				err = tx.Get(&tagID, "SELECT (id) FROM content.image_tags WHERE description = $1;",
					tag)
				if err != nil {
					log.Printf("CREATE NEW TAG: %s\n", tag)
					rows, err := tx.Query(`INSERT INTO content.image_tags (description)
											VALUES ($1) RETURNING id;`, tag)
					if err != nil {
						log.Println(err)
						return err
					}
					for rows.Next() {
						rows.Scan(&tagID)
					}

					rows.Close()
				}
				_, err = tx.Exec(`INSERT INTO content.image_tag_bridge (image_id, tag_id)
										VALUES ($1, $2);`, image.Id, tagID)
				if err != nil {
					log.Println(err)
					return err
				}
			}
		} else {
			_, err = tx.Exec(fmt.Sprintf(`UPDATE content.image_metadata SET %s = $1 WHERE image_id = $2;`, key), val, image.Id)
			if err != nil {
				log.Println(err)
				log.Printf(`UPDATE content.image_metadata SET %s = $1 WHERE image_id = $2;`, key)
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
		log.Printf("UPDATE content.users SET %s = %s WHERE id = %d", key, val, user.Id)
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
