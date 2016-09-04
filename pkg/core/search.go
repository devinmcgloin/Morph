package core

// import (
// 	"net/http"
//
// 	"github.com/sprioc/composer/pkg/model"
// 	"github.com/sprioc/composer/pkg/qmgo"
// 	"github.com/sprioc/composer/pkg/rsp"
// 	"github.com/sprioc/composer/pkg/store"
// )
//
// func Search(user model.User, query qmgo.ImageSearch) ([]*model.Image, rsp.Response) {
// 	if !query.Valid() {
// 		return []*model.Image{}, rsp.Response{Code: http.StatusBadRequest}
// 	}
//
// 	var images []*model.Image
//
// 	err := store.SearchImages(query, &images)
// 	if err != nil {
// 		return []*model.Image{}, rsp.Response{Code: http.StatusInternalServerError}
// 	}
//
// 	return images, rsp.Response{Code: http.StatusOK}
// }
