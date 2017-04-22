package sql

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/sprioc/composer/pkg/model"
)

func Permissions(userRef uint32, permission model.Permission, item uint32) (bool, error) {
	var valid int
	var err error
	var stmt *sqlx.Stmt

	switch permission {
	case model.CanEdit:
		stmt, err = db.Preparex("SELECT count(*) FROM permissions.can_edit WHERE user_id = ? AND o_id = ?;")
	case model.CanView:
		stmt, err = db.Preparex("SELECT count(*) FROM permissions.can_view WHERE user_id = ? AND o_id = ?;")
	case model.CanDelete:
		stmt, err = db.Preparex("SELECT count(*) FROM permissions.can_delete WHERE user_id = ? AND o_id = ?;")
	}
	if err != nil {
		log.Println(err)
		return false, err
	}
	err = stmt.Get(&valid, userRef, item)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return valid == 1, nil

}

func AddPermissions(item uint32, permission model.Permission, userRef uint32) error {
	var err error
	var stmt *sqlx.Stmt

	switch permission {
	case model.CanEdit:
		stmt, err = db.Preparex("INSERT INTO permissions.can_edit(user_id, o_id) VALUES(?, ?);")
	case model.CanView:
		stmt, err = db.Preparex("INSERT INTO permissions.can_view(user_id, o_id) VALUES(?, ?);")
	case model.CanDelete:
		stmt, err = db.Preparex("INSERT INTO permissions.can_delete(user_id, o_id) VALUES(?, ?);")
	}
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(userRef, item)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// IsAdmin checks if the given user has admin privileges
func IsAdmin(id uint32) (bool, error) {
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

// GetLogin returns the salt, password, email and username for a given user.
func GetLogin(ref uint32) (map[string]string, error) {

	rows, err := db.Query("SELECT (id, username, salt, password, email) FROM content.users WHERE username = ? LIMIT 1;", ref)
	if err != nil {
		log.Print(err)
		return map[string]string{}, err
	}
	defer rows.Close()
	var userInfo = map[string]string{}

	for rows.Next() {
		if err := rows.Scan(&userInfo); err != nil {
			log.Print(err)
			return map[string]string{}, err
		}
	}
	if err := rows.Err(); err != nil {
		log.Print(err)
		return map[string]string{}, err
	}

	return userInfo, nil
}
