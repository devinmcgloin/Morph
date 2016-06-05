package dbase

import (
	"database/sql"
	"os"

	"github.com/devinmcgloin/morph/src/schema"
	_ "github.com/go-sql-driver/mysql" // want sql drivers to init, work with the database/sql package.

	"log"
)

// GetDB returns a reference to a sql.DB object. It's best to keep these long lived.
func GetDB() *sql.DB {

	// Create the database handle, confirm driver is
	db, err := sql.Open("mysql", Getenv("JAWSDB_MARIA_URL", "root:@/morph"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetPage(imdID string, db *sql.DB) schema.ImgPage {

	var page schema.ImgPage
	rows, err := db.Query("select ", args)

	return page
}

func GetCollectionPage(collectionTag string, db *sql.DB) schema.ImgCollection {

	var collectionPage schema.ImgCollection

	return collectionPage
}

// Getenv is the same as os.Getenv but allows for secondary key.
func Getenv(key string, opt string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return opt
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
