package modification

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func deleteUser(db *sqlx.DB, id int64) error {
	var imageIds []int64

	err := db.Select(&imageIds, "SELECT id FROM content.images WHERE owner = $1", id)
	if err != nil {
		log.Println(err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	for _, id := range imageIds {
		tx.Exec("DELETE FROM content.image_metadata WHERE image_id = $1", id)
		tx.Exec("DELETE FROM content.image_label_bridge WHERE image_id = $1", id)
		tx.Exec("DELETE FROM content.image_tag_bridge WHERE image_id = $1", id)
		tx.Exec("DELETE FROM content.user_favorites WHERE image_id = $1", id)
		tx.Exec("DELETE FROM content.image_color_bridge WHERE image_id = $1", id)
		tx.Exec("DELETE FROM content.image_landmark_bridge WHERE image_id = $1", id)
		tx.Exec("DELETE FROM content.image_geo WHERE image_id = $1", id)
		tx.Exec("DELETE FROM content.images WHERE id = $1", id)
	}
	tx.Exec("DELETE FROM content.users WHERE id = $1", id)
	return tx.Commit()
}

// DeleteImage removes all keys for the given image, as well as removing it from
// the owner. In the future it will also handle favorites and collections.
func deleteImage(db *sqlx.DB, id int64) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Print(err)
		return err
	}
	log.Printf("Deleting image_id = %d\n", id)

	tx.Exec("DELETE FROM content.image_metadata WHERE image_id = $1", id)
	tx.Exec("DELETE FROM content.image_label_bridge WHERE image_id = $1", id)
	tx.Exec("DELETE FROM content.image_tag_bridge WHERE image_id = $1", id)
	tx.Exec("DELETE FROM content.user_favorites WHERE image_id = $1", id)
	tx.Exec("DELETE FROM content.image_color_bridge WHERE image_id = $1", id)
	tx.Exec("DELETE FROM content.image_landmark_bridge WHERE image_id = $1", id)
	tx.Exec("DELETE FROM content.image_geo WHERE image_id = $1", id)
	tx.Exec("DELETE FROM content.images WHERE id = $1", id)
	err = tx.Commit()
	return err
}
