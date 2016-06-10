package main

import (
	"github.com/devinmcgloin/morph/src/env"
	_ "github.com/go-sql-driver/mysql" // want sql drivers to init, work with the database/sql package.
	"github.com/jmoiron/sqlx"

	"log"
)

var imageSchema = `
CREATE TABLE IF NOT EXISTS images
  (
     i_id            INT UNSIGNED NOT NULL auto_increment,
     i_title         TEXT DEFAULT NULL,
     i_desc          TEXT DEFAULT NULL,
     i_aperture      INT UNSIGNED DEFAULT NULL,
     i_exposure_time INT UNSIGNED DEFAULT NULL,
     i_focal_length  INT UNSIGNED DEFAULT NULL,
     i_iso           INT UNSIGNED DEFAULT NULL,
     i_orientation   TEXT DEFAULT NULL,
     i_camera_body   TEXT DEFAULT NULL,
     i_lens          TEXT DEFAULT NULL,
     i_tag_1         TEXT DEFAULT NULL,
     i_tag_2         TEXT DEFAULT NULL,
     i_tag_3         TEXT DEFAULT NULL,
     a_id            TEXT DEFAULT NULL,
     i_capture_time  DATETIME DEFAULT NULL,
     i_publish_time  DATETIME DEFAULT CURRENT_TIMESTAMP,
     i_direction     FLOAT DEFAULT NULL,
     u_id            INT UNSIGNED NOT NULL,
     l_id            INT UNSIGNED DEFAULT NULL,
     PRIMARY KEY(i_id)
  );
		`

var userSchema = `
CREATE TABLE IF NOT EXISTS users
  (
     u_id         INT UNSIGNED NOT NULL auto_increment,
     u_username   TEXT DEFAULT NULL,
     u_email      TEXT DEFAULT NULL,
     u_first_name TEXT DEFAULT NULL,
     u_last_name  TEXT DEFAULT NULL,
		 u_bio        TEXT DEFAULT NULL,
		 u_avatar_url TEXT DEFAULT NULL,
		 l_id         INT UNSIGNED DEFAULT NULL,
     PRIMARY KEY(u_id)
  );
`

var sourceSchema = `
CREATE TABLE IF NOT EXISTS sources
  (
     s_id         INT UNSIGNED NOT NULL auto_increment,
     i_id         INT UNSIGNED NOT NULL,
     s_url        TEXT NOT NULL,
     s_resolution INT UNSIGNED DEFAULT NULL,
     s_width      INT UNSIGNED DEFAULT NULL,
     s_height 	  INT UNSIGNED DEFAULT NULL,
     s_size       TEXT DEFAULT NULL,
     s_file_type  TEXT DEFAULT NULL,
     PRIMARY KEY(s_id)
  );
`

var locationSchema = `
CREATE TABLE IF NOT EXISTS locations
  (
     l_id   INT UNSIGNED NOT NULL auto_increment,
     l_desc TEXT DEFAULT NULL,
     l_lat  FLOAT DEFAULT NULL,
     l_lon  FLOAT DEFAULT NULL,
     PRIMARY KEY(l_id)
  );
`

var albumSchema = `
CREATE TABLE IF NOT EXISTS albums
  (
     a_id       INT UNSIGNED NOT NULL auto_increment,
     u_id       INT UNSIGNED NOT NULL,
     a_desc     TEXT DEFAULT NULL,
     a_title    TEXT DEFAULT NULL,
     a_viewtype TEXT DEFAULT NULL,
     PRIMARY KEY(a_id)
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

	db.MustExec(albumSchema)

	log.Println("Tables Initialized")

}
