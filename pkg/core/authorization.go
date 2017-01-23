package core

import (
	"strings"

	"github.com/sprioc/composer/pkg/model"
)

var validOps = []string{"$set", "$unset"}

func validTarget(ref model.Ref, target string) bool {
	imgTarget := []string{"tags", "featured", "aperature",
		"exposure_time", "focal_length", "iso",
		"make", "model", "lens_make", "lens_model"}
	userTarget := []string{"name", "bio", "personal_site_link"}
	albEvntColTarget := []string{"desc", "title", "view_type"}

	switch ref.Collection {
	case model.Images:
		return in(target, imgTarget)
	case model.Users:
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
