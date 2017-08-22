package create

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"errors"

	"fmt"

	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/request"
	"github.com/devinmcgloin/fokal/pkg/retrieval"
	"github.com/jmoiron/sqlx"
)

var validEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Anything but special characters and spaces.
var validUsername = regexp.MustCompile(`^[^\<\>\!\{\}\[\]\!\@\#\$\%\^\&\*\(\)\.\ ]{3,16}$`)

var letters = regexp.MustCompile("^[a-zA-Z]+$")

func validateUser(db *sqlx.DB, userData *request.CreateUserRequest) error {
	userData.Username = strings.ToLower(userData.Username)

	if !validUsername.MatchString(userData.Username) {
		return handler.StatusError{Err: errors.New("Invalid username"), Code: http.StatusBadRequest}
	}

	userData.Email = strings.Trim(strings.ToLower(userData.Email), " ")

	if !validEmail.MatchString(userData.Email) {
		return handler.StatusError{Err: errors.New("Invalid email"), Code: http.StatusBadRequest}
	}

	if !(validPassword(userData.Password) || validPassPhrase(userData.Password)) {
		return handler.StatusError{Err: errors.New("Invalid password"), Code: http.StatusBadRequest}
	}

	userRef := model.Ref{Collection: model.Users, Shortcode: userData.Username}

	exists, err := retrieval.ExistsUser(db, userRef.Shortcode)
	if err != nil {
		return handler.StatusError{Code: http.StatusInternalServerError,
			Err: errors.New(http.StatusText(http.StatusInternalServerError))}
	}
	if exists {
		return handler.StatusError{Err: errors.New("Username already exist"), Code: http.StatusConflict}
	}

	exists, err = retrieval.ExistsEmail(db, userData.Email)
	if err != nil {
		return handler.StatusError{Code: http.StatusInternalServerError}
	}
	if exists {
		return handler.StatusError{Err: errors.New("Email already exist"), Code: http.StatusConflict}
	}

	return nil
}

func validPassword(password string) bool {
	var hasUpper bool
	var hasLower bool
	var hasSpecial bool
	var hasNumber bool

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		default:
			return false
		}
	}

	return hasLower && hasUpper && hasNumber && hasSpecial && len(password) > 8
}

func validPassPhrase(passphrase string) bool {
	sections := strings.Split(passphrase, "-")

	for _, sect := range sections {
		if !letters.MatchString(sect) {
			return false
		} else if len(sect) < 5 {
			return false
		}
	}
	return len(sections) >= 3
}

// CreateImage stores the image data in the database under the given user.
// Currently does not set the metadata or db interal state.
func commitImage(db *sqlx.DB, image model.Image) error {
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

	_, err = tx.Exec(`
	INSERT INTO content.image_geo (image_id, loc, dir)
	VALUES ($1, GeomFromEWKB($2), $3);
	`, image.Id, image.Metadata.Location, image.Metadata.ImageDirection)
	if err != nil {
		log.Println(err)
		return err
	}

	// Adding landmarks
	var landmarkID int64
	for _, landmark := range image.Landmarks {
		err := tx.Get(&landmarkID, "SELECT id FROM content.landmarks WHERE description = $1", landmark.Description)
		if err != nil {

			err = tx.Get(&landmarkID, `
			INSERT INTO content.landmarks(description, location)
			VALUES($1, GeomFromEWKB($2)) RETURNING id;`, landmark.Description,
				landmark.Location)
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
			l, a, b := color.SRGB.CIELAB()
			lab := fmt.Sprintf("(%f, %f, %f)", l, a, b)

			err = tx.Get(&colorID, `
			INSERT INTO content.colors (red, green, blue, hue, saturation, val, shade, color, cielab)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::cube) RETURNING id;`, color.SRGB.R, color.SRGB.G, color.SRGB.B,
				color.HSV.H, color.HSV.S, color.HSV.V, color.Shade, color.ColorName, lab)
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
			err = tx.Get(&labelID, `INSERT INTO content.labels (description) VALUES($1) RETURNING id;`,
				label.Description)
			if err != nil {
				log.Println(err)
				return err
			}
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

func CommitUser(db *sqlx.DB, username, email, name string) error {
	var uID int64
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	rows, err := tx.Query(`
	INSERT INTO content.users(username, email, name)
	VALUES($1, $2, $3) RETURNING id;`,
		username, email, name)
	if err != nil {
		log.Println(err)
		return err
	}

	for rows.Next() {
		err = rows.Scan(&uID)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	rows.Close()

	_, err = tx.Exec(`
	INSERT INTO permissions.can_edit(user_id, o_id, type) VALUES ($1, $1, 'user');`, uID)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.Exec(`
	INSERT INTO permissions.can_delete(user_id, o_id, type) VALUES ($1, $1, 'user');`, uID)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.Exec(`
	INSERT INTO permissions.can_view(user_id, o_id, type) VALUES (-1, $1, 'user');`, uID)
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
