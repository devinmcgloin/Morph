package sql

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/sprioc/composer/pkg/model"
)

func Permissions(userRef int64, permission model.Permission, item int64) (bool, error) {
	var valid int
	var err error
	var stmt *sqlx.Stmt

	switch permission {
	case model.CanEdit:
		stmt, err = db.Preparex("SELECT count(*) FROM permissions.can_edit WHERE user_id = $1 AND o_id = $1;")
	case model.CanView:
		stmt, err = db.Preparex("SELECT count(*) FROM permissions.can_view WHERE user_id = $1 AND o_id = $1;")
	case model.CanDelete:
		stmt, err = db.Preparex("SELECT count(*) FROM permissions.can_delete WHERE user_id = $1 AND o_id = $1;")
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

func AddPermissions(userRef int64, permission model.Permission, item int64) error {
	var err error
	var stmt *sqlx.Stmt

	switch permission {
	case model.CanEdit:
		stmt, err = db.Preparex("INSERT INTO permissions.can_edit(user_id, o_id) VALUES($1, $2);")
	case model.CanView:
		stmt, err = db.Preparex("INSERT INTO permissions.can_view(user_id, o_id) VALUES($1, $2);")
	case model.CanDelete:
		stmt, err = db.Preparex("INSERT INTO permissions.can_delete(user_id, o_id) VALUES($1, $2);")
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
func IsAdmin(id int64) (bool, error) {
	rows, err := db.Query("SELECT count(*) FROM content.users WHERE id = $1 AND admin = TRUE", id)
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
func GetLogin(ref string) (map[string]interface{}, error) {
	userInfo := make(map[string]interface{})
	rows, err := db.Query("SELECT id, username, salt, password, email FROM content.users WHERE username = $1 LIMIT 1;", ref)
	if err != nil {
		log.Print(err)
		return userInfo, err
	}
	defer rows.Close()
	var id int64
	var username string
	var salt string
	var password string
	var email string
	for rows.Next() {
		if err := rows.Scan(&id, &username, &salt, &password, &email); err != nil {
			log.Print(err)
			return userInfo, err
		}
	}
	if err := rows.Err(); err != nil {
		log.Print(err)
		return userInfo, err
	}
	userInfo["id"] = id
	userInfo["username"] = username
	userInfo["salt"] = salt
	userInfo["password"] = password
	userInfo["email"] = email

	return userInfo, nil
}
