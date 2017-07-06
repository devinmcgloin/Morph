package social

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func Favorite(db *sqlx.DB, uID, iID int64) error {
	stmt, err := db.Preparex("INSERT INTO content.user_favorites (user_id, image_id) VALUES ($1, $2)")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(db, stmt, uID, iID)
}

func UnFavorite(db *sqlx.DB, uID, iID int64) error {
	stmt, err := db.Preparex("DELETE FROM content.user_favorites WHERE user_id = $1 AND image_id = $2")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(db, stmt, uID, iID)
}

func Follow(db *sqlx.DB, idA, idB int64) error {
	stmt, err := db.Preparex("INSERT INTO content.user_follows (user_id, user_follows) VALUES ($1, $2)")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(db, stmt, idA, idB)
}

func UnFollow(db *sqlx.DB, idA, idB int64) error {
	stmt, err := db.Preparex("DELETE FROM content.user_follows WHERE user_id = $1 AND user_follows = $2)")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(db, stmt, idA, idB)
}

func AddTag(db *sqlx.DB, iID, tagID int64) error {
	stmt, err := db.Preparex("INSERT INTO content.image_tag_bridge (image_id, tag_id) VALUES ($1, $2)")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(db, stmt, iID, tagID)
}

func RemoveTag(db *sqlx.DB, iID, tagID int64) error {
	stmt, err := db.Preparex("DELETE FROM content.image_tag_bridge WHERE image_id = $1 AND tag_id = $2")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(db, stmt, iID, tagID)
}

func modify(db *sqlx.DB, stmt *sqlx.Stmt, idA int64, idB int64) error {
	_, err := stmt.Exec(idA, idB)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
