package core

import "strings"

var validOps = []string{"$set", "$unset"}

const (
	invalid = iota - 1
	images
	users
	collections
	albums
	events
)

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
