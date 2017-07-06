package modification

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func deleteUser(db *sqlx.DB, id int64) error {
	var image_ids []int64

	err := db.Select(&image_ids, "SELECT id FROM content.images WHERE owner = ?", id)
	if err != nil {
		log.Println(err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	tx.Exec("DELETE FROM content.users WHERE id = ?", id)
	for _, img := range image_ids {
		tx.Exec("DELETE FROM content.images WHERE id = ?", img)
		tx.Exec("DELETE FROM content.image_metadata WHERE image_id = ?", img)
		tx.Exec("DELETE FROM content.image_label_bridge WHERE image_id = ?", img)
		tx.Exec("DELETE FROM content.image_tag_bridge WHERE image_id = ?", img)
		tx.Exec("DELETE FROM content.user_favorites WHERE image_id = ?", img)
	}
	return tx.Commit()
}

// DeleteImage removes all keys for the given image, as well as removing it from
// the owner. In the future it will also handle favorites and collections.
func deleteImage(db *sqlx.DB, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	tx.Exec("DELETE FROM content.images WHERE id = ?", id)
	tx.Exec("DELETE FROM content.image_metadata WHERE image_id = ?", id)
	tx.Exec("DELETE FROM content.image_label_bridge WHERE image_id = ?", id)
	tx.Exec("DELETE FROM content.image_tag_bridge WHERE image_id = ?", id)
	tx.Exec("DELETE FROM content.user_favorites WHERE image_id = ?", id)
	return tx.Commit()
}
