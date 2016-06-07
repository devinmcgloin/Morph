package dbase

import (
	"strings"

	"github.com/devinmcgloin/morph/src/env"
	"github.com/devinmcgloin/morph/src/schema"
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

func GetImg(iID string) (schema.Img, error) {

	var img schema.Img

	err := db.Get(&img, "SELECT * FROM images WHERE i_id = ?", iID)
	if err != nil {
		return schema.Img{}, nil
	}

	log.Printf("Image: %v", img)

	return img, nil
}

func AddImg(img schema.Img) error {

	_, err := db.NamedExec(`
		INSERT INTO images (i_id, i_title, i_desc, i_url, i_category, i_fstop, i_shutter_speed, i_fov, i_iso, i_publish_date)
			VALUES (:id, :title, :desc, :url, :category, :fstop, :shutter, :fov, :iso, :publish_date)`,
		map[string]interface{}{
			"id":           img.IID,
			"title":        img.Title,
			"desc":         img.Desc,
			"url":          img.URL,
			"category":     img.Category,
			"fstop":        img.FStop,
			"shutter":      img.ShutterSpeed,
			"fov":          img.FOV,
			"iso":          img.ISO,
			"publish_date": img.PublishDate,
		})
	if err != nil {
		return err
	}
	return nil

}

func GetCategory(collectionTag string) (schema.ImgCollection, error) {
	var collectionPage schema.ImgCollection

	var images []schema.Img

	err := db.Select(&images, "SELECT * FROM images WHERE i_category = ?", collectionTag)
	if err != nil {
		return schema.ImgCollection{}, err
	}

	collectionPage.Images = images
	collectionPage.NumImg = len(collectionPage.Images)
	collectionPage.Title = collectionTag
	return collectionPage, nil
}

func GetAllImgs() (schema.ImgCollection, error) {
	var collectionPage schema.ImgCollection

	var images []schema.Img

	err := db.Select(&images, "SELECT * FROM images ORDER BY i_publish_date DESC")
	if err != nil {
		return schema.ImgCollection{}, err
	}

	collectionPage.Images = images
	collectionPage.NumImg = len(collectionPage.Images)
	collectionPage.Title = "Morph"
	return collectionPage, nil
}

func generateQuestionMarks(num int) string {
	var qs []string
	for i := 0; i < num; i++ {
		qs = append(qs, "?")
	}
	return strings.Join(qs, ", ")
}
