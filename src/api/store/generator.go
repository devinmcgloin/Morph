package store

import (
	"github.com/devinmcgloin/sprioc/src/generator"
	"gopkg.in/mgo.v2/bson"
)

func (ds *MgoStore) GetNewImageID() bson.ObjectId {
	id := bson.ObjectId(generator.RandString(12))
	if ds.ExistsImageID(id) {
		id = bson.ObjectId(generator.RandString(12))
	}
	return id
}

func (ds *MgoStore) GetNewAlbumID() bson.ObjectId {
	id := bson.ObjectId(generator.RandString(12))
	if ds.ExistsAlbumID(id) {
		id = bson.ObjectId(generator.RandString(12))
	}
	return id
}

func (ds *MgoStore) GetNewEventID() bson.ObjectId {
	id := bson.ObjectId(generator.RandString(12))
	if ds.ExistsEventID(id) {
		id = bson.ObjectId(generator.RandString(12))
	}
	return id
}

func (ds *MgoStore) GetNewCollectionID() bson.ObjectId {
	id := bson.ObjectId(generator.RandString(12))
	if ds.ExistsCollectionID(id) {
		id = bson.ObjectId(generator.RandString(12))
	}
	return id
}
