package main

import (
	"log"

	"github.com/devinmcgloin/morph/src/env"
	_ "github.com/go-sql-driver/mysql" // want sql drivers to init, work with the database/sql package.
	"github.com/jmoiron/sqlx"
)

func main() {
	log.Println("Connecting to db...")

	// Create the database handle, confirm driver is
	db, err := sqlx.Open("mysql", env.Getenv("DB_URL", "root:@/morph"))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Clearing db tables...")

	db.Exec("DROP TABLE users")

	db.Exec("DROP TABLE images")

	db.Exec("DROP TABLE sources")

	db.Exec("DROP TABLE locations")

	db.Exec("DROP TABLE albums")

	log.Println("Tables Cleared")

}
