package core

import (
	"net/http"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/redis"

	"github.com/sprioc/composer/pkg/rsp"
)

func ModifySecure(user model.User, target model.Ref, changes [][]string) rsp.Response {

	// checking if the user has permission to modify the item
	valid, err := redis.Permissions(user.GetRef(), model.CanEdit, target)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError, Message: "Unable to retrieve user permissions."}
	}
	if !valid {
		return rsp.Response{Code: http.StatusForbidden, Message: "User does not have permission to edit item."}
	}

	// checking if modification is valid.
	return rsp.Response{Code: http.StatusNotImplemented}
}
