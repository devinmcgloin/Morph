package main

import (
	_ "github.com/go-sql-driver/mysql" // want sql drivers to init, work with the database/sql package.
	"github.com/jmoiron/sqlx"

	"log"

	"github.com/devinmcgloin/morph/src/env"
)

var imageSchema = `
CREATE TABLE IF NOT EXISTS images
  (
     i_id            INT NOT NULL auto_increment,
     i_title         TEXT DEFAULT NULL,
     i_desc          TEXT DEFAULT NULL,
     i_aperture      INT DEFAULT NULL,
     i_exposure_time INT DEFAULT NULL,
     i_focal_length  INT DEFAULT NULL,
     i_iso           INT DEFAULT NULL,
     i_orientation   TEXT DEFAULT NULL,
     i_camera_body   TEXT DEFAULT NULL,
     i_lens          TEXT DEFAULT NULL,
     i_tag_1         TEXT DEFAULT NULL,
     i_tag_2         TEXT DEFAULT NULL,
     i_tag_3         TEXT DEFAULT NULL,
     i_album         TEXT DEFAULT NULL,
     i_capture_time  DATETIME DEFAULT NULL,
     i_publish_time  DATETIME DEFAULT CURRENT_TIMESTAMP,
     i_direction     FLOAT DEFAULT NULL,
     u_id            INT NOT NULL,
     l_id            INT DEFAULT NULL,
     PRIMARY KEY(i_id)
  );
		`

var userSchema = `
CREATE TABLE IF NOT EXISTS users
  (
     u_id         INT NOT NULL auto_increment,
     u_username   TEXT DEFAULT NULL,
     u_email      TEXT DEFAULT NULL,
     u_first_name TEXT DEFAULT NULL,
     u_last_name  TEXT DEFAULT NULL,
     PRIMARY KEY(u_id)
  );
`

var sourceSchema = `
CREATE TABLE IF NOT EXISTS source
  (
     s_id         INT NOT NULL auto_increment,
     i_id         INT DEFAULT NULL,
     i_url        TEXT DEFAULT NULL,
     u_resolution INT DEFAULT NULL,
     u_width      INT DEFAULT NULL,
	 s_length 	  INT DEFAULT NULL,
	 s_size       INT DEFAULT 0,
	 s_file_type  INT DEFAULT NULL,
     PRIMARY KEY(s_id)
  );
`

var locationSchema = `
CREATE TABLE IF NOT EXISTS locations
  (
     l_id   INT NOT NULL auto_increment,
     l_desc TEXT DEFAULT NULL,
     l_lat  FLOAT DEFAULT NULL,
     l_lon  FLOAT DEFAULT NULL,
     PRIMARY KEY(l_id)
  );
`

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

	log.Println("Initializing db tables...")

	db.MustExec(imageSchema)

	db.MustExec(userSchema)

	db.MustExec(sourceSchema)

	db.MustExec(locationSchema)

	log.Println("Tables Initialized")

}
