package sql

import (
	"log"
)

// This can be refactored using the sqlx get operation

// ExistsImage checks if the given id exists in the database
func ExistsImage(shortcode string) (bool, error) {
	rows, err := db.Query("SELECT count(*) FROM content.images WHERE shortcode = $1;", shortcode)
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
func ExistsUser(username string) (bool, error) {
	count := 0
	err := db.Get(&count, "SELECT count(*) FROM content.users WHERE username = $1;", username)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return count == 1, nil
}

// ExistsEmail checks if there is a user record with the given email
func ExistsEmail(email string) (bool, error) {
	rows, err := db.Query("SELECT count(*) FROM content.users WHERE email = $1;", email)
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
