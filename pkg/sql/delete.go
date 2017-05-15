package sql

import (
	"log"
)

// DeleteUser removes all references to a given user, as well as the images
// they have uploaded. In the future it will also remove following relationships
// and other graph relationships.
func DeleteUser(username string) error {
	var id int64
	var image_ids []int64

	err := db.Select(&id, "SELECT id FROM content.users WHERE username = ?", username)
	if err != nil {
		log.Println(err)
		return err
	}

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
func DeleteImage(shortcode string) error {
	var id int64

	err := db.Select(&id, "SELECT id FROM content.images WHERE shortcode = ?", shortcode)
	if err != nil {
		log.Println(err)
		return err
	}

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
