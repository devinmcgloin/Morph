package sql

import (
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // Importing the package for its side effect.
)

func Configure(postgresURL string) {
	var err error

	db, err = sqlx.Open("postgres", postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

var db *sqlx.DB
