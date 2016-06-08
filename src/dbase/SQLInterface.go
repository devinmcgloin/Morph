package dbase

import (
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
			INSERT INTO images (
				i_id,
				i_title,
				i_desc,
				i_aperture,
				i_exposure_time,
				i_focal_length,
				i_iso,
				i_orientation,
				i_camera_body,
				i_lens,
				i_tag_1,
				i_tag_2,
				i_tag_3,
				i_album,
				i_capture_time,
				i_publish_time,
				i_lat,
				i_lon,
				i_direction,
				i_loc,
				i_user)
		VALUES (
			:i_id,
			:i_title,
			:i_desc,
			:i_aperture,
			:i_exposure_time,
			:i_focal_length,
			:i_iso,
			:i_orientation,
			:i_camera_body,
			:i_lens,
			:i_tag_1,
			:i_tag_2,
			:i_tag_3,
			:i_album,
			:i_capture_time,
			:i_publish_time,
			:i_lat,
			:i_lon,
			:i_direction,
			:i_loc,
			:i_user)`,
		map[string]interface{}{
			":i_id":            img.ID,
			":i_title":         img.Title,
			":i_desc":          img.Desc,
			":i_aperture":      img.Aperture,
			":i_exposure_time": img.ExposureTime,
			":i_focal_length":  img.FocalLength,
			":i_iso":           img.ISO,
			":i_orientation":   img.Orientation,
			":i_camera_body":   img.CameraBody,
			":i_lens":          img.Lens,
			":i_tag_1":         img.TagOne,
			":i_tag_2":         img.TagTwo,
			":i_tag_3":         img.TagThree,
			":i_album":         img.Album,
			":i_capture_time":  img.CaptureTime,
			":i_publish_time":  img.PublishTime,
			":i_direction":     img.ImgDirection,
			":l_loc":           img.Location,
			":u_user":          img.User,
		})
	log.Printf("DB Error Type = %s", err.Error())
	if err != nil {
		return err
	}
	return nil
}

func GetAlbum(albumTag string) (schema.ImgCollection, error) {
	var collectionPage schema.ImgCollection

	var images []schema.Img

	err := db.Select(&images, "SELECT * FROM images WHERE i_album = ?", albumTag)
	if err != nil {
		return schema.ImgCollection{}, err
	}

	collectionPage.Images = images
	collectionPage.NumImg = len(collectionPage.Images)
	collectionPage.Title = albumTag
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
