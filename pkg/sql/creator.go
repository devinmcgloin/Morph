package sql

import (
	"log"

	postgis "github.com/cridenour/go-postgis"
	"github.com/sprioc/composer/pkg/model"
)

// CreateImage stores the image data in the database under the given user.
// Currently does not set the metadata or db interal state.
func CreateImage(image model.Image) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}
	var id int64
	rows, err := tx.NamedQuery(`
	INSERT INTO content.images(user_id, shortcode)
	VALUES(:user_id, :shortcode) RETURNING id;`,
		image)
	if err != nil {
		log.Println(err)
		return err
	}

	for rows.Next() {
		rows.Scan(&id)
	}

	rows.Close()

	image.Id = id

	_, err = tx.NamedExec(`
	INSERT INTO content.image_metadata(image_id, aperture, exposure_time, 
	focal_length, iso, make, model, lens_make, lens_model, pixel_xd, 
	pixel_yd, capture_time) VALUES (:id, :metadata.aperture, :metadata.exposure_time,
	:metadata.focal_length, :metadata.iso, :metadata.make, :metadata.model,
	:metadata.lens_make, :metadata.lens_model, :metadata.pixel_xd,
	:metadata.pixel_yd, :metadata.capture_time);
	`, image)
	if err != nil {
		log.Println(err)
		return err
	}

	point := postgis.PointS{SRID: 4326,
		X: float64(image.Metadata.Location.Coordinates[0]),
		Y: float64(image.Metadata.Location.Coordinates[1])}
	_, err = tx.Exec(`
	INSERT INTO content.image_geo (image_id, loc, dir) 
	VALUES ($1, GeomFromEWKB($2), $3);
	`, image.Id, point, image.Metadata.ImageDirection)
	if err != nil {
		log.Println(err)
		return err
	}

	// Adding landmarks
	var landmarkID int64
	for _, landmark := range image.Landmarks {
		err := tx.Get(&landmarkID, "SELECT id FROM content.landmarks WHERE desc = $1", landmark.Description)
		if err != nil {
			log.Println(err)
			return err
		}
		if landmarkID == 0 {
			point := postgis.PointS{SRID: 4326,
				X: float64(landmark.Location.Coordinates[0]),
				Y: float64(landmark.Location.Coordinates[1])}

			err = tx.Get(&landmarkID, `
			INSERT INTO content.landmarks(desc, location)
			VALUES($1, GeomFromEWKB($2)) RETURNING id;`, landmark.Description,
				point)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		_, err = tx.Exec(`
			INSERT INTO content.image_landmark_bridge(image_id, landmark_id, score)
			VALUES ($1, $2, $3)`, image.Id, landmarkID, landmark.Score)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// Adding colors
	var colorID int64
	for _, color := range image.Colors {
		err := tx.Get(&colorID, "SELECT id FROM content.colors "+
			"WHERE red = $1 AND green = $2 AND blue = $3", color.SRGB.R, color.SRGB.G, color.SRGB.B)
		if err != nil {
			log.Println(err)
			return err
		}
		if colorID == 0 {
			err = tx.Get(&colorID, `
			INSERT INTO content.colors (red, green, blue, hue, saturation, val, shade, color)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`, color.SRGB.R, color.SRGB.G, color.SRGB.B,
				color.HSV.H, color.HSV.S, color.HSV.V, color.Shade, color.ColorName)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		_, err = tx.Exec(`
			INSERT INTO content.image_color_bridge(image_id, color_id, pixel_fraction, score)
			VALUES ($1, $2, $3, $4)`, image.Id, colorID, color.PixelFraction, color.Score)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// Adding Labels
	var labelID int64
	for _, label := range image.Labels {
		err := tx.Get(&labelID, ` SELECT id FROM content.labels WHERE description = $1`,
			label.Description)
		if err != nil {
			log.Println(err)
			return err
		}

		if labelID == 0 {
			err = tx.Get(&labelID, `INSERT INTO content.labels (description) VALUES($1) RETURNING id;`,
				label.Description)
		}

		_, err = tx.Exec(`
			INSERT INTO content.image_label_bridge(image_id, label_id, score)
			VALUES ($1, $2, $3)`, image.Id, labelID, label.Score)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// Permissions
	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_edit(user_id, o_id, type) VALUES (:user_id, :id, 'image');
	`, image)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_delete(user_id, o_id, type) VALUES (:user_id, :id, 'image');
	`, image)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_view(user_id, o_id, type) VALUES (-1, :id, 'image');
	`, image)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			return err
		}
		return err
	}
	return nil
}

func CreateUser(user model.User) error {
	log.Println(user)

	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	rows, err := tx.NamedQuery(`
	INSERT INTO content.users(username, email, password, salt)
	VALUES(:username, :email, :password, :salt) RETURNING id;`,
		user)
	if err != nil {
		log.Println(err)
		return err
	}

	for rows.Next() {
		err = rows.Scan(&user.Id)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	rows.Close()

	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_edit(user_id, o_id, type) VALUES (:id, :id, 'user');`, user)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_delete(user_id, o_id, type) VALUES (:id, :id, 'user');`, user)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.NamedExec(`
	INSERT INTO permissions.can_view(user_id, o_id, type) VALUES (-1, :id, 'user');`, user)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			return err
		}
		return err
	}

	return nil
}
