package core

import (
	"net/http"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/redis"
	"github.com/sprioc/composer/pkg/rsp"
)

func GetUser(ref model.Ref) (model.User, rsp.Response) {
	if ref.Collection != model.Users {
		return model.User{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}

	user, err := redis.GetUser(ref, false)
	if err != nil {
		return model.User{}, rsp.Response{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return user, rsp.Response{Code: http.StatusOK}
}

func GetImage(ref model.Ref) (model.Image, rsp.Response) {
	if ref.Collection != model.Images {
		return model.Image{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}

	image, err := redis.GetImage(ref)
	if err != nil {
		return model.Image{}, rsp.Response{Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return image, rsp.Response{Code: http.StatusOK}
}
