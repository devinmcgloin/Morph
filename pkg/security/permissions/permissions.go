package permissions

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Permission string

const (
	CanEdit   = Permission("can_edit")
	CanDelete = Permission("can_delete")
	CanView   = Permission("can_view")
)

func Valid(db *sqlx.DB, userRef int64, permission Permission, item int64) (bool, error) {
	var valid int
	var err error
	var stmt *sqlx.Stmt

	switch permission {
	case CanEdit:
		stmt, err = db.Preparex("SELECT count(*) FROM permissions.can_edit WHERE user_id = $1 AND o_id = $2;")
	case CanView:
		stmt, err = db.Preparex("SELECT count(*) FROM permissions.can_view WHERE (user_id = $1 OR user_id = -1) AND o_id = $2;")
	case CanDelete:
		stmt, err = db.Preparex("SELECT count(*) FROM permissions.can_delete WHERE user_id = $1 AND o_id = $2;")
	}
	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Printf("usr: %d, permission: %v, item: %d ", userRef, permission, item)
	err = stmt.Get(&valid, userRef, item)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if valid == 0 {
		admin, err := IsAdmin(db, userRef)
		if err != nil {
			log.Println(err)
			return false, err
		}
		if admin {
			return admin, nil
		}
	}
	log.Printf("valid: %d\n", valid)
	return valid == 1, nil

}

func Add(db *sqlx.DB, userRef int64, permission Permission, item int64) error {
	var err error
	var stmt *sqlx.Stmt

	switch permission {
	case CanEdit:
		stmt, err = db.Preparex("INSERT INTO permissions.can_edit(user_id, o_id) VALUES($1, $2);")
	case CanView:
		stmt, err = db.Preparex("INSERT INTO permissions.can_view(user_id, o_id) VALUES($1, $2);")
	case CanDelete:
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
func IsAdmin(db *sqlx.DB, id int64) (bool, error) {
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
