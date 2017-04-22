package sql

import (
	"log"

	"github.com/sprioc/composer/pkg/model"
)

// CreateImage stores the image data in the database under the given user.
// Currently does not set the metadata or db interal state.
func CreateImage(image model.Image) error {
	sc, err := generateSC(model.Images)
	if err != nil {
		log.Println(err)
		return err
	}
	image.Shortcode = sc

	_, err = db.NamedExec(`
	INSERT INTO content.images(owner, shortcode)
	VALUES(:owner, :shortcode);`,
		image)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func CreateUser(user model.User) error {
	_, err := db.NamedExec(`
	INSERT INTO content.users(username, email, password, salt)
	VALUES(:username, :email, :password, :salt);`,
		user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
