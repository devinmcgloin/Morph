package sql

import (
	"log"

	"github.com/sprioc/composer/pkg/model"
)

func Permissions(userRef model.UserReference, permission model.Permission, item model.ImageReference) (bool, error) {
}

func AddPermissions(item model.ImageReference, permission model.Permission, userRef model.UserReference) error {
}

// IsAdmin checks if the given user has admin privileges
func IsAdmin(id model.UserReference) (bool, error) {
	rows, err := db.Query("SELECT count(*) FROM content.users WHERE id = ? AND admin = TRUE", id)
	if err != nil {
		log.Print(err)
		return false, err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Print(err)
			return false, err
		}
	}
	if err := rows.Err(); err != nil {
		log.Print(err)
		return false, err
	}

	return count == 1, nil
}

func GetLogin(ref model.UserReference) (map[string]string, error) {

}
