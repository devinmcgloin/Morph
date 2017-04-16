package sql

import (
	"errors"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/refs"
)

func GetUser(u model.Ref, priv bool) (model.User, error) {

}

func GetImage(i model.Ref) (model.Image, error) {
	img := model.Image{}
	err := db.Select(img, "SELECT * FROM content.images WHERE id = ?", i.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

}
