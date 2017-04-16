package sql

import "log"

// ExistsImage checks if the given id exists in the database
func ExistsImage(id string) (bool, error) {
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
func ExistsUser(id string) (bool, error) {
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
