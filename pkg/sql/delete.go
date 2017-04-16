package sql

import (
	"github.com/sprioc/composer/pkg/model"
	"log"
)

// DeleteUser removes all references to a given user, as well as the images
// they have uploaded. In the future it will also remove following relationships
// and other graph relationships.
func DeleteUser(ref model.Ref) error {
	var image_ids []string
	err := db.Select(image_ids, "SELECT id FROM content.images WHERE owner = ?", ref.Id)
	if err != nil {
		log.Println(err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	tx.Exec("DELETE FROM content.users WHERE id = ?", ref.Id)
	for _, img := range image_ids {
		tx.Exec("DELETE FROM content.images WHERE id = ?", img)
		tx.Exec("DELETE FROM content.image_tags WHERE image_id = ?", img)
		tx.Exec("DELETE FROM content.image_labels WHERE image_id = ?", img)
	}
	return tx.Commit()
}

// DeleteImage removes all keys for the given image, as well as removing it from
// the owner. In the future it will also handle favorites and collections.
func DeleteImage(ref model.Ref) error {
	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}
	tx.Exec("DELETE FROM content.images WHERE id = ?", ref.Id)
	tx.Exec("DELETE FROM content.image_tags WHERE image_id = ?", ref.Id)
	tx.Exec("DELETE FROM content.image_labels WHERE image_id = ?", ref.Id)
	return tx.Commit()
}
