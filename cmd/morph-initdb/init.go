package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // want sql drivers to init, work with the database/sql package.

	"log"

	"github.com/devinmcgloin/morph/src/env"
)

func main() {
	log.Println("Connecting to db...")

	// Create the database handle, confirm driver is
	db, err := sql.Open("mysql", env.Getenv("DB_URL", "root:@/morph"))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing db tables...")

	createPhotosTable := `
	CREATE TABLE IF NOT EXISTS images
	  (
	     i_id            INT NOT NULL auto_increment,
	     i_title         TEXT(30) DEFAULT NULL,
	     i_desc          TEXT(255) DEFAULT NULL,
	     i_url           TEXT DEFAULT NULL,
	     i_fstop         INT DEFAULT NULL,
	     i_iso           INT DEFAULT NULL,
	     i_fov           INT DEFAULT NULL,
	     i_shutter_speed INT DEFAULT NULL,
	     i_category      TEXT DEFAULT NULL,
	     i_publish_date  DATETIME DEFAULT CURRENT_TIMESTAMP,
	     PRIMARY KEY(i_id)
	  );
`
	executeSQL(db, createPhotosTable)

	createConfigTable := `
	CREATE TABLE IF NOT EXISTS config
	  (
	     conf_lock           CHAR(1) NOT NULL DEFAULT 'X',
	     conf_author         TEXT(30) DEFAULT NULL,
	     conf_author_twitter TEXT(50) DEFAULT NULL,
	     conf_desc           TEXT DEFAULT NULL,
	     conf_keywords       TEXT DEFAULT NULL,
	     CONSTRAINT pk_config PRIMARY KEY (conf_lock),
	     CONSTRAINT ck_config_locked CHECK (conf_lock='X')
	  );
	`
	executeSQL(db, createConfigTable)

	log.Println("Table creation completed.")

	db.Close()

}

func executeSQL(db *sql.DB, sqlCreateTable string) {
	stmt, err := db.Prepare(sqlCreateTable)
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)

}
