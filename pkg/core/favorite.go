package core

import (
	"net/http"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/redis"
	"github.com/sprioc/composer/pkg/rsp"
)

// TODO would like to tell users the request failed due to the target not existsing.

func Favorite(user model.Ref, target model.Ref) rsp.Response {

	if !user.Valid(model.Users) {
		return rsp.Response{Message: "Only users can favorite things.", Code: http.StatusBadRequest}
	}

	if user.ShortCode == target.ShortCode {
		return rsp.Response{Message: "User cannot favorite themselves.", Code: http.StatusBadRequest}
	}

	err := redis.LinkItems(user, redis.Favorite, target, false)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}

}

func UnFavorite(user model.Ref, target model.Ref) rsp.Response {

	if !user.Valid(model.Users) {
		return rsp.Response{Message: "Only users can favorite things.", Code: http.StatusBadRequest}
	}

	if user.ShortCode == target.ShortCode {
		return rsp.Response{Message: "User cannot favorite themselves.", Code: http.StatusBadRequest}
	}

	err := redis.LinkItems(user, redis.Favorite, target, true)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}
