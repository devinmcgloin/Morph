package main

import (
	"log"

	"github.com/devinmcgloin/morph/src/env"
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

	db.MustExec("DROP TABLE users")

	db.MustExec("DROP TABLE images")

	db.MustExec("DROP TABLE sources")

	db.MustExec("DROP TABLE locations")

	db.MustExec("DROP TABLE albums")

	log.Println("Tables Cleared")

}
