package sql

import (
	"log"
)

// DeleteUser removes all references to a given user, as well as the images
// they have uploaded. In the future it will also remove following relationships
// and other graph relationships.
func DeleteUser(id int64) error {
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
		//tx.Exec("DELETE FROM content.image_tags WHERE image_id = ?", img)
		//tx.Exec("DELETE FROM content.image_labels WHERE image_id = ?", img)
	}
	return tx.Commit()
}

// DeleteImage removes all keys for the given image, as well as removing it from
// the owner. In the future it will also handle favorites and collections.
func DeleteImage(id int64) error {
	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}
	tx.Exec("DELETE FROM content.images WHERE id = ?", id)
	//tx.Exec("DELETE FROM content.image_tags WHERE image_id = ?", ref)
	//tx.Exec("DELETE FROM content.image_labels WHERE image_id = ?", ref)
	return tx.Commit()
}
