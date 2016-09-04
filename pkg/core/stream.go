package core

// import (
// 	"net/http"
//
// 	"github.com/sprioc/composer/pkg/model"
// 	"github.com/sprioc/composer/pkg/rsp"
// )
//
// // GetStream finds the images that are recently posted in the collections and
// // users that this user follows.
// func GetStream(user model.User) ([]*model.Image, rsp.Response) {
//
// 	var images []*model.Image
//
// 	for _, ref := range user.Followes {
// 		switch ref.Collection {
// 		case "collections":
// 			followImg, resp := GetCollectionImages(ref)
// 			if resp.Ok() {
// 				images = append(images, followImg...)
// 			}
// 		case "users":
// 			followImg, resp := GetUserImages(ref)
// 			if resp.Ok() {
// 				images = append(images, followImg...)
// 			}
// 		}
// 	}
//
// 	return images, rsp.Response{Code: http.StatusOK}
// }
