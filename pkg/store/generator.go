package store

import "github.com/sprioc/conductor/pkg/generator"

func GetNewImageShortCode() string {
	id := generator.RandString(12)
	if ExistsImageID(id) {
		id = generator.RandString(12)
	}
	return id
}

func GetNewAlbumShortCode() string {
	id := generator.RandString(12)
	if ExistsAlbumID(id) {
		id = generator.RandString(12)
	}
	return id
}

func GetNewEventShortCode() string {
	id := generator.RandString(12)
	if ExistsEventID(id) {
		id = generator.RandString(12)
	}
	return id
}

func GetNewCollectionShortCode() string {
	id := generator.RandString(12)
	if ExistsCollectionID(id) {
		id = generator.RandString(12)
	}
	return id
}

func GetNewUserShortCode() string {
	id := generator.RandString(12)
	if ExistsUserID(id) {
		id = generator.RandString(12)
	}
	return id
}
