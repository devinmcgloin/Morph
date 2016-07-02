package core

import (
	"log"
	"net/http"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func GetStream(user model.User) ([]*model.Image, rsp.Response) {

	var images []*model.Image

	for _, ref := range user.Followes {
		log.Println(ref)
		switch ref.Collection {
		case "collections":
			followImg, resp := GetCollectionImages(ref)
			if resp.Ok() {
				log.Println(followImg)
				images = append(images, followImg...)
			}
		case "users":
			followImg, resp := GetUserImages(ref)
			if resp.Ok() {
				images = append(images, followImg...)
				log.Println(followImg)
			}
		}
	}

	log.Println(len(images))
	return images, rsp.Response{Code: http.StatusOK}

}
