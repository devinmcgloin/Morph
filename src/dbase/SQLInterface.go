package dbase

import (
	"database/sql"

	"github.com/devinmcgloin/morph/src/env"
	"github.com/devinmcgloin/morph/src/schema"
	_ "github.com/go-sql-driver/mysql" // want sql drivers to init, work with the database/sql package.

	"log"
)

var DB *sql.DB

// SetDB returns a reference to a sql.DB object. It's best to keep these long lived.
func SetDB() *sql.DB {

	// Create the database handle, confirm driver is
	db, err := sql.Open("mysql", env.Getenv("DB_URL", "root:@/morph"))
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	return db
}

func GetImg(pID string, db *sql.DB) schema.Img {

	var page schema.Img

	rows, err := db.Query(
		`
			SELECT p_title,
			       p_desc,
			       p_url,
			       p_fstop,
			       p_iso,
			       p_fov,
			       p_shutter_speed,
			       p_category
			FROM   photos
			WHERE  p_id = ?
			`, pID)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&page.Title, &page.Desc, &page.URL, &page.PhotoMeta.FStop,
			&page.PhotoMeta.ISO, &page.PhotoMeta.FOV, &page.PhotoMeta.ShutterSpeed, &page.Category)

		if err != nil {
			log.Fatal(err)
		}
		log.Println(page)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return page
}

func AddImg(img schema.Img, db *sql.DB) {

	stmt, err := db.Prepare(
		`INSERT INTO photos
            (p_id,
						 p_title,
             p_desc,
             p_url,
             p_fstop,
             p_iso,
             p_fov,
             p_shutter_speed,
             p_category,
					   p_publish_date)
		VALUES      (?, ?, ?, ?, ?,
			           ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Fatalf("db.Prepare failed %s", err)
	}

	res, err := stmt.Exec("NULL", img.Title, img.Desc, img.URL, img.PhotoMeta.FStop, img.PhotoMeta.ISO,
		img.PhotoMeta.FOV, img.PhotoMeta.ShutterSpeed, img.Category, "NULL")
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
	stmt.Close()
}

func getCollection(collectionTag string, db *sql.DB) schema.ImgCollection {
	var collectionPage schema.ImgCollection

	return collectionPage
}

func GetAllImgs(db *sql.DB) schema.ImgCollection {
	var collectionPage schema.ImgCollection

	var page schema.Img

	rows, err := db.Query(
		`
			SELECT p_title,
			       p_desc,
			       p_url,
			       p_fstop,
			       p_iso,
			       p_fov,
			       p_shutter_speed,
			       p_category
			FROM   photos
			`)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&page.Title, &page.Desc, &page.URL, &page.PhotoMeta.FStop,
			&page.PhotoMeta.ISO, &page.PhotoMeta.FOV, &page.PhotoMeta.ShutterSpeed, &page.Category)

		if err != nil {
			log.Fatal(err)
		}
		collectionPage.Images = append(collectionPage.Images, page)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	collectionPage.NumImg = len(collectionPage.Images)
	collectionPage.Title = "Home"
	return collectionPage
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
