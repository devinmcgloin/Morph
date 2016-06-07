package dbase

import (
	"database/sql"
	"strings"

	"github.com/devinmcgloin/morph/src/env"
	"github.com/devinmcgloin/morph/src/schema"

	_ "github.com/go-sql-driver/mysql" // Drives we have to import to use database/sql

	"log"
)

var db *sql.DB

// SetDB returns a reference to a sql.DB object. It's best to keep these long lived.
func SetDB() error {
	log.Printf("DB_URL = %s", env.Getenv("DB_URL", "root:@/morph"))

	var err error
	// Create the database handle, confirm driver is
	db, err = sql.Open("mysql", env.Getenv("DB_URL", "root:@/morph")+"?parseTime=true")
	if err != nil {
		log.Print(err)
		return err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to db, ping failed.")
	}
	return nil
}

func GetImg(iID string) (schema.Img, error) {

	var img schema.Img

	err := db.QueryRow("select * from images where i_id = ?", iID).Scan(
		&img.IID, &img.Title, &img.Desc, &img.URL, &img.FStop,
		&img.ISO, &img.FOV, &img.ShutterSpeed, &img.Category, &img.PublishDate,
	)
	if err != nil {
		return schema.Img{}, err
	}

	log.Printf("Image: %v", img)

	return img, nil
}

func AddImg(img schema.Img) error {

	stmt, err := db.Prepare("INSERT INTO images(" + generateQuestionMarks(10) +
		") VALUES(" + generateQuestionMarks(10) + ")")
	if err != nil {
		log.Fatal(err)
		return err
	}

	res, err := stmt.Exec("i_id", "i_title", "i_desc", "i_url", "i_fstop",
		"i_iso", "i_fov", "i_shutter_speed", "i_category", "i_publish_date",
		img.IID, img.Title, img.Desc, img.URL, img.FStop, img.ISO, img.FOV,
		img.ShutterSpeed, img.Category, img.PublishDate)

	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(res)

	return nil

}

func GetCategory(collectionTag string) (schema.ImgCollection, error) {
	var collectionPage schema.ImgCollection

	var images []schema.Img

	rows, err := db.Query("SELECT "+generateQuestionMarks(10)+" FROM images WHERE i_category = ?", args)

	collectionPage.Images = images
	collectionPage.NumImg = len(collectionPage.Images)
	collectionPage.Title = "Morph"
	return collectionPage, nil
}

func GetAllImgs() (schema.ImgCollection, error) {
	var collectionPage schema.ImgCollection
	//
	// var images []schema.Img
	//
	// sql := "SELECT i_id, i_title, i_desc, i_url, i_category, i_fstop, i_shutter_speed, i_fov, i_iso, i_publish_date name FROM images ORDER BY i_publish_date DESC"
	//
	// q := db.NewQuery(sql)
	// err := q.All(&images)
	// if err != nil {
	// 	return schema.ImgCollection{}, err
	// }
	//
	// collectionPage.Images = images
	// collectionPage.NumImg = len(collectionPage.Images)
	// collectionPage.Title = "Morph"
	return collectionPage, nil
}

func generateQuestionMarks(num int) string {
	var qs []string
	for i := 0; i < num; i++ {
		qs = append(qs, "?")
	}
	return strings.Join(qs, ", ")
}


func getAllDBFields(schema struct) string {
	
}
