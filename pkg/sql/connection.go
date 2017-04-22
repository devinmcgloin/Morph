package sql

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // Importing the package for its side effect.
)

func init() {
	connURL = os.Getenv("POST_URL")

	if connURL == "" {
		log.Fatal("POST_URL not set")
	}

	db, err := sqlx.Open("postgres", connURL)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connected Succesfully to Postgres DB\n")
}

var connURL string
var db *sqlx.DB
