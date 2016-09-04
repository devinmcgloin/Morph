package core

//
// import (
// 	"net/http"
//
// 	"github.com/sprioc/composer/pkg/model"
// 	"github.com/sprioc/composer/pkg/rsp"
// 	"github.com/sprioc/composer/pkg/store"
// 	"gopkg.in/mgo.v2/bson"
// )
//
// func FeatureImage(user model.User, image model.Ref) rsp.Response {
//
// 	if !user.Admin {
// 		return rsp.Response{Code: http.StatusForbidden, Message: "Only admins can feature images"}
// 	}
//
// 	err := store.Modify(image, bson.M{"$set": bson.M{"featured": true}})
// 	if err != nil {
// 		return rsp.Response{Code: http.StatusInternalServerError}
// 	}
//
// 	return rsp.Response{Code: http.StatusAccepted}
// }
//
// func UnFeatureImage(user model.User, image model.Ref) rsp.Response {
// 	if !user.Admin {
// 		return rsp.Response{Code: http.StatusForbidden, Message: "Only admins can feature images"}
// 	}
//
// 	err := store.Modify(image, bson.M{"$set": bson.M{"featured": false}})
// 	if err != nil {
// 		return rsp.Response{Code: http.StatusInternalServerError}
// 	}
//
// 	return rsp.Response{Code: http.StatusAccepted}
// }
