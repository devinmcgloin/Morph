package core

import (
	"net/http"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/redis"
	"github.com/sprioc/composer/pkg/rsp"
)

func FeatureImage(user model.Ref, image model.Ref) rsp.Response {
	admin, err := redis.IsAdmin(user)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	if !admin {
		return rsp.Response{Code: http.StatusForbidden, Message: "Only admins can feature images"}
	}

	err = redis.Set(image.GetRString(model.Featured), true)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func UnFeatureImage(user model.Ref, image model.Ref) rsp.Response {
	admin, err := redis.IsAdmin(user)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	if !admin {
		return rsp.Response{Code: http.StatusForbidden, Message: "Only admins can feature images"}
	}

	err = redis.Set(image.GetRString(model.Featured), false)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}
