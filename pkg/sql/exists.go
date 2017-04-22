package sql

import (
	"log"

	"github.com/sprioc/composer/pkg/model"
)

// This can be refactored using the sqlx get operation

// ExistsImage checks if the given id exists in the database
func ExistsImage(id model.ImageReference) (bool, error) {
	rows, err := db.Query("SELECT count(*) FROM content.images WHERE id = ?", id)
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

// ExistsUser checks if the given id exists in the database
func ExistsUser(id model.UserReference) (bool, error) {
	rows, err := db.Query("SELECT count(*) FROM content.users WHERE id = ?", id)
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

// ExistsEmail checks if there is a user record with the given email
func ExistsEmail(email string) (bool, error) {
	rows, err := db.Query("SELECT count(*) FROM content.users WHERE email = ?", email)
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
