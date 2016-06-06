package dbase

import (
	"time"

	"github.com/devinmcgloin/morph/src/env"
	"github.com/devinmcgloin/morph/src/schema"
	"github.com/go-ozzo/ozzo-dbx"

	_ "github.com/go-sql-driver/mysql"

	"log"
)

var db *dbx.DB

// SetDB returns a reference to a sql.DB object. It's best to keep these long lived.
func SetDB() error {
	log.Printf("DB_URL = %s", env.Getenv("DB_URL", "root:@/morph"))

	var err error
	// Create the database handle, confirm driver is
	db, err = dbx.Open("mysql", env.Getenv("DB_URL", "root:@/morph"))
	if err != nil {
		log.Print(err)
		return err
	}
	err = db.DB().Ping()
	if err != nil {
		log.Fatal("Error connecting to db, ping failed.")
	}
	log.Printf("Database successfully launched with %s driver.", db.DriverName())
	return nil
}

func GetImg(iID string) (schema.Img, error) {

	var img schema.Img

	sql := "SELECT i_id, i_title, i_desc, i_url, i_category, i_fstop, i_shutter_speed, i_fov, i_iso, i_publish_date name FROM images WHERE i_id={:iID}"

	q := db.NewQuery(sql)
	err := q.Bind(dbx.Params{"iID": iID}).One(&img)
	if err != nil {
		return schema.Img{}, err
	}

	log.Printf("Image: %v", img)

	return img, nil
}

func AddImg(img schema.Img) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Insert("images", dbx.Params{
		"i_title":         img.Title,
		"i_desc":          img.Desc,
		"i_url":           img.URL,
		"i_category":      img.Category,
		"i_fstop":         img.FStop,
		"i_shutter_speed": img.ShutterSpeed,
		"i_fov":           img.FOV,
		"i_iso":           img.ISO,
		"i_publish_date":  time.Now(),
	}).Execute()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func GetCategory(collectionTag string) (schema.ImgCollection, error) {
	var collectionPage schema.ImgCollection

	var images []schema.Img

	sql := "SELECT i_id, i_title, i_desc, i_url, i_category, i_fstop, i_shutter_speed, i_fov, i_iso, i_publish_date name FROM images WHERE i_category={:category}"

	q := db.NewQuery(sql)
	err := q.Bind(dbx.Params{"category": collectionTag}).All(&images)
	if err != nil {
		return schema.ImgCollection{}, err
	}

	collectionPage.Images = images
	collectionPage.NumImg = len(collectionPage.Images)
	collectionPage.Title = "Morph"
	return collectionPage, nil
}

func GetAllImgs() (schema.ImgCollection, error) {
	var collectionPage schema.ImgCollection

	var images []schema.Img

	sql := "SELECT i_id, i_title, i_desc, i_url, i_category, i_fstop, i_shutter_speed, i_fov, i_iso, i_publish_date name FROM images ORDER BY i_publish_date DESC"

	q := db.NewQuery(sql)
	err := q.All(&images)
	if err != nil {
		return schema.ImgCollection{}, err
	}

	log.Printf("Length of all Images = %d", len(images))

	collectionPage.Images = images
	collectionPage.NumImg = len(collectionPage.Images)
	collectionPage.Title = "Morph"
	return collectionPage, nil
}
