package core

import (
	"reflect"

	"github.com/sprioc/composer/pkg/model"
)

// REVIEW this code is narly and slow, not sure how well the db will maintain
// internal consistency. Mongo doesn't seem to give many gurantees.

// func DeleteImage(requestFrom model.Ref, imageRef model.Ref) rsp.Response {
//
// }
//
// func DeleteUser(requestFrom model.User, user model.Ref) rsp.Response {
//
// }

func inRef(item model.Ref, collection []model.Ref) bool {
	for _, x := range collection {
		if reflect.DeepEqual(x, item) {
			return true
		}
	}
	return false
}
