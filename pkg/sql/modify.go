package sql

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func Favorite(username, shortcode string) error {
	stmt := db.Preparex("INSERT INTO content.user_favorites (user_id, image_id) VALUES ($1, $2)")

	uID := GetUserID(username)
	iID := GetImageID(shortcode)
	return modify(stmt, uID, iID)
}

func UnFavorite(username, shortcode string) error {
	stmt := db.Preparex("DELETE FROM content.user_favorites WHERE user_id = $1 AND image_id = $2")
	uID := GetUserID(username)
	iID := GetImageID(shortcode)

	return modify(stmt, uID, iID)
}

func Follow(userA, userB string) error {
	stmt := db.Preparex("INSERT INTO content.user_follows (user_id, user_follows) VALUES ($1, $2)")
	idA := GetUserID(userA)
	idB := GetUserID(userB)

	return modify(stmt, idA, idB)
}

func UnFollow(userA, userB string) error {
	stmt := db.Preparex("DELETE FROM content.user_follows WHERE user_id = $1 AND user_follows = $2)")
	idA := GetUserID(userA)
	idB := GetUserID(userB)

	return modify(stmt, idA, idB)
}

func AddTag(shortcode, tag string) error {
	stmt := db.Preparex("INSERT INTO content.image_tag_bridge (image_id, tag_id) VALUES ($1, $2)")
	iID := GetImageID(shortcode)
	var tagID int64

	err := db.Get(&tagID, "SELECT id FROM content.image_tags WHERE description = $1", tag)
	if err != nil {
		log.Println(err)
		return err
	}

	return modify(stmt, iID, tagID)
}

func RemoveTag(shortcode, tag string) error {
	stmt := db.Preparex("DELETE FROM content.image_tag_bridge WHERE image_id = $1 AND tag_id = $1")
	iID := GetImageID(shortcode)
	var tagID int64

	err := db.Get(&tagID, "SELECT id FROM content.image_tags WHERE description = $1", tag)
	if err != nil {
		log.Println(err)
		return err
	}

	return modify(stmt, iID, tagID)
}

func Feature(shortcode string) error {
	iID := GetImageID(shortcode)

	err := db.Exec(`
	UPDATE content.images
		SET featured = TRUE
	WHERE id = $1`, iID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func UnFeature(shortcode string) error {
	iID := GetImageID(shortcode)

	err := db.Exec(`
	UPDATE content.images
		SET featured = FALSE
	WHERE id = $1`, iID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func modify(stmt *sqlx.Stmt, idA int64, idB int64) error {
	err = stmt.Exec(idA, idB)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
