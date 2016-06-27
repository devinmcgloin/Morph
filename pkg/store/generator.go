package store

import "github.com/sprioc/sprioc-core/pkg/generator"

func (ds *MgoStore) GetNewImageShortCode() string {
	id := generator.RandString(12)
	if ds.ExistsImageID(id) {
		id = generator.RandString(12)
	}
	return id
}

func (ds *MgoStore) GetNewAlbumShortCode() string {
	id := generator.RandString(12)
	if ds.ExistsAlbumID(id) {
		id = generator.RandString(12)
	}
	return id
}

func (ds *MgoStore) GetNewEventShortCode() string {
	id := generator.RandString(12)
	if ds.ExistsEventID(id) {
		id = generator.RandString(12)
	}
	return id
}

func (ds *MgoStore) GetNewCollectionShortCode() string {
	id := generator.RandString(12)
	if ds.ExistsCollectionID(id) {
		id = generator.RandString(12)
	}
	return id
}

func (ds *MgoStore) GetNewUserShortCode() string {
	id := generator.RandString(12)
	if ds.ExistsUserID(id) {
		id = generator.RandString(12)
	}
	return id
}
