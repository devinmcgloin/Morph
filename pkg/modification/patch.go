package modification

import (
	"fmt"
	"log"

	postgis "github.com/cridenour/go-postgis"
	"github.com/fokal/fokal-core/pkg/model"
	"github.com/jmoiron/sqlx"
)

func commitImagePatch(db *sqlx.DB, image model.Ref, req map[string]interface{}) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	for key, val := range req {
		if key == "tags" {
			_, err = tx.Exec("DELETE FROM content.image_tag_bridge WHERE image_id = $1;", image.Id)
			if err != nil {
				log.Println(err)
				return err
			}

			tags := val.([]string)
			var tagID int
			for _, tag := range tags {
				err = tx.Get(&tagID, "SELECT (id) FROM content.image_tags WHERE description = LOWER($1);",
					tag)
				if err != nil {
					log.Printf("CREATE NEW TAG: %s\n", tag)
					rows, err := tx.Query(`INSERT INTO content.image_tags (description)
											VALUES (LOWER($1)) RETURNING id;`, tag)
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
		} else if key == "geo" {
			loc := val.(map[string]interface{})
			p := postgis.PointS{
				SRID: 4326,
				X:    loc["Longitude"].(float64),
				Y:    loc["Latitude"].(float64),
			}
			desc := loc["Description"].(string)
			log.Println(image, p, desc)
			_, err = tx.Exec(`
			INSERT INTO content.image_geo (image_id, loc, description)
			VALUES ($1, GeomFromEWKB($2), $3)
				ON CONFLICT (image_id) DO UPDATE
					SET loc = excluded.loc,
						description = excluded.description`,
				image.Id, p, desc)
			if err != nil {
				log.Println(err)
				return err
			}
		} else {
			_, err = tx.Exec(fmt.Sprintf(`UPDATE content.image_metadata SET %s = $1 WHERE image_id = $2;`, key), val, image.Id)
			log.Printf("Setting %s to %v\n", key, val)
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

func commitUserPatch(db *sqlx.DB, user model.Ref, changes map[string]interface{}) error {
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
