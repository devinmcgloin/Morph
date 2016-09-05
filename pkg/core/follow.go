package core

import (
	"net/http"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/redis"
	"github.com/sprioc/composer/pkg/rsp"
)

// TODO would like to tell users the request failed due to the target not existsing.

func Follow(user model.Ref, target model.Ref) rsp.Response {

	if target.Valid(model.Collections, model.Images, model.Users) || user.Valid(model.Users) {
		return rsp.Response{Code: http.StatusBadRequest, Message: "Invalid References"}
	}

	if user.ShortCode == target.ShortCode {
		return rsp.Response{Message: "User cannot follow themselves.", Code: http.StatusBadRequest}
	}

	err := redis.LinkItems(user, redis.Follow, target, false)
	if err != nil {
		return rsp.Response{Message: "Error while adding relation", Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func UnFollow(user model.Ref, target model.Ref) rsp.Response {

	if target.Valid(model.Collections, model.Images, model.Users) || user.Valid(model.Users) {
		return rsp.Response{Code: http.StatusBadRequest, Message: "Invalid References"}
	}

	if user.ShortCode == target.ShortCode {
		return rsp.Response{Message: "User cannot follow themselves.", Code: http.StatusBadRequest}
	}

	err := redis.LinkItems(user, redis.Follow, target, true)
	if err != nil {
		return rsp.Response{Message: "Error while adding relation", Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}
