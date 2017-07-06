package modification

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func Feature(db *sqlx.DB, iID int64) error {
	_, err := db.Exec(`
	UPDATE content.images
		SET featured = TRUE
	WHERE id = $1`, iID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func UnFeature(db *sqlx.DB, iID int64) error {
	_, err := db.Exec(`
	UPDATE content.images
		SET featured = FALSE
	WHERE id = $1`, iID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
