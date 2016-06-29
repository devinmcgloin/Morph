package core

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/rsp"
	"gopkg.in/mgo.v2/bson"
)

var validOps = []string{"$set", "$unset"}

const (
	invalid = iota - 1
	images
	users
	collections
	albums
	events
)

// VerifyChanges checks operations, permissions and targets to see if an
// operation is valid. Should not be used for internal modifications.
func VerifyChanges(user model.User, target model.DBRef, changes bson.M) rsp.Response {
	// admins may need to edit other desc or posts
	if user.Admin {
		return rsp.Response{Code: http.StatusOK}
	}

	var targetType int

	targetType, resp := Authorized(user, target)
	if !resp.Ok() {
		return resp
	}

	for key, val := range changes {
		if !in(key, validOps) {
			return rsp.Response{Message: fmt.Sprintf("Invalid Operation: %s", key), Code: http.StatusBadRequest}
		}

		m, ok := val.(map[string]interface{})
		if !ok {
			return rsp.Response{Message: "Invalid changefile", Code: http.StatusBadRequest}
		}

		for tar := range m {
			if !validTarget(targetType, tar) {
				return rsp.Response{Message: fmt.Sprintf("Invalid Target: %s", tar), Code: http.StatusBadRequest}
			}
		}
	}
	return rsp.Response{Code: http.StatusOK}
}

// Authorized checks only if the account can modify the given target.
func Authorized(user model.User, target model.DBRef) (int, rsp.Response) {
	switch {
	case target.Collection == "images":
		if inRef(target, user.Images) {
			return images, rsp.Response{Code: http.StatusOK}
		}
		break
	case target.Collection == "users":
		if target.Shortcode == user.ShortCode {
			return users, rsp.Response{Code: http.StatusOK}
		}
		break
	case target.Collection == "collections":
		if inRef(target, user.Collections) {
			return collections, rsp.Response{Code: http.StatusOK}
		}
		break
	default:
		return invalid, rsp.Response{Message: fmt.Sprintf("Invalid target type %s", target), Code: http.StatusBadRequest}
	}
	return invalid, rsp.Response{Message: fmt.Sprintf("User has invalid credentials"), Code: http.StatusUnauthorized}
}

func validTarget(kind int, target string) bool {
	imgTarget := []string{"tags", "featured", "metadata.aperature",
		"metadata.exposure_time", "metadata.focal_length", "metadata.iso",
		"metadata.make", "metadata.model", "metadata.lens_make", "metadata.lens_model"}
	userTarget := []string{"name", "bio", "personal_site_link"}
	albEvntColTarget := []string{"desc", "title", "view_type"}

	switch kind {
	case images:
		return in(target, imgTarget)
	case users:
		return in(target, userTarget)
	default:
		return in(target, albEvntColTarget)
	}
}

func in(contentType string, opts []string) bool {
	for _, opt := range opts {
		if strings.Compare(contentType, opt) == 0 {
			return true
		}
	}
	return false
}
