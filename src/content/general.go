package content

import (
	"github.com/devinmcgloin/morph/src/env"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql" // Drives we have to import to use database/sql

	"log"
)

var db *sqlx.DB

// SetDB returns a reference to a sql.DB object. It's best to keep these long lived.
func SetDB() error {
	log.Printf("DB_URL = %s", env.Getenv("DB_URL", "root:@/morph"))

	var err error
	// Create the database handle, confirm driver is
	db, err = sqlx.Connect("mysql", env.Getenv("DB_URL", "root:@/morph")+"?parseTime=true")
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
