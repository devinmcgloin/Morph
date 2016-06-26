package authorization

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sprioc/sprioc-core/pkg/model"
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
func VerifyChanges(user model.User, target interface{}, changes bson.M, internal bool) error {
	// admins may need to edit other desc or posts
	if user.Admin {
		return nil
	}

	var targetType int

	targetType, hasAuth := Authorized(user, target)
	if hasAuth != nil {
		return hasAuth
	}

	// If we reach this point it means the user has authorization to change the
	// document, as its an internal modification we can bypass verifying the changes
	if internal {
		return nil
	}

	for key, val := range changes {
		if !in(key, validOps) {
			return fmt.Errorf("Invalid Operation: %s", key)
		}

		m, ok := val.(map[string]interface{})
		if !ok {
			return errors.New("Invalid changefile")
		}

		for tar := range m {
			if !validTarget(targetType, tar) {
				return fmt.Errorf("Invalid Target: %s", key)
			}
		}
	}
	return nil
}

// Authorized checks only if the account can modify the given target.
func Authorized(user model.User, target interface{}) (int, error) {
	switch t := target.(type) {
	case model.Image:
		if strings.Compare(user.ShortCode, target.(model.Image).User.Shortcode) == 0 {
			return images, nil
		}
		break
	case model.User:
		if strings.Compare(user.ShortCode, target.(model.User).ShortCode) == 0 {
			return users, nil
		}
		break
	case model.Collection:
		if strings.Compare(user.ShortCode, target.(model.Collection).Curator.Shortcode) == 0 {
			return collections, nil
		}
		break
	case model.Album:
		if strings.Compare(user.ShortCode, target.(model.Album).User.Shortcode) == 0 {
			return albums, nil
		}
		break
	case model.Event:
		return events, fmt.Errorf("User has invalid credentials")
	default:
		return invalid, fmt.Errorf("Invalid target type %s", t)
	}
	return invalid, fmt.Errorf("User has invalid credentials")
}

func validTarget(kind int, target string) bool {
	imgTarget := []string{"tags", "featured", "metadata.aperature",
		"metadata.exposure_time", "metadata.focal_length", "metadata.iso",
		"metadata.make", "metadata.model", "metadata.lens_make", "metadata.lens_model"}
	userTarget := []string{"email", "name", "bio", "url"}
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
