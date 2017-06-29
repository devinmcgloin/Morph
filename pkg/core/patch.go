package core

import (
	"net/http"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/rsp"
	"github.com/sprioc/composer/pkg/sql"
)

// PatchImage aplies the requesting changes to the image.
func PatchImage(user model.Ref, image model.Ref, request map[string]interface{}) rsp.Response {
	resp := permission(user, model.CanEdit, image)
	if !resp.Ok() {
		return resp
	}

	valid := make(map[string]interface{})
	dest := [10]string{"tags", "aperature", "exposure_time", "focal_length",
		"iso", "make", "model", "lens_make", "lens_model", "capture_time"}

	for _, loc := range dest {
		change, ok := request[loc]
		if ok {
			valid[loc] = change
		}
	}

	tags, notPresent := valid["tags"]
	_, ok := tags.([]string)
	if !ok && !notPresent {
		return rsp.Response{Code: http.StatusBadRequest, Message: "Tags must be an array of new values"}
	}

	err := sql.PatchImage(image, valid)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}

func PatchUser(user model.Ref, target model.Ref, request map[string]interface{}) rsp.Response {
	resp := permission(user, model.CanEdit, user)
	if !resp.Ok() {
		return resp
	}

	valid := make(map[string]interface{})
	dest := [3]string{"bio", "url", "name"}

	for _, loc := range dest {
		change, ok := request[loc]
		if ok {
			valid[loc] = change
		}
	}

	err := sql.PatchUser(target, valid)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}
