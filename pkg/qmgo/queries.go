package qmgo

// Sets of relationships
import (
	"time"

	"github.com/sprioc/conductor/pkg/model"
	"github.com/sprioc/conductor/pkg/refs"
	"github.com/sprioc/geojson"
)

type ImageSearch struct {
	PublishTime map[Ord]time.Time `json:"publish_time" bson:"publish_time,omitempty"`
	Featured    bool              `json:"featured" bson:"featured,omitempty"`
	Downloads   map[Ord]int       `json:"downloads" bson:"downloads,omitempty"`
	Owner       string            `json:"owner" bson:"-"`
	OwnerExtern model.DBRef       `json:"-" bson:"owner,omitempty"`
	MetaData    struct {
		CaptureTime map[Ord]time.Time     `json:"capture_time" bson:"capture_time,omitempty"`
		Location    map[Geo]geojson.Point `json:"location" bson:"location,omitempty"`
	} `json:"metadata" bson:"metadata,omitempty"`
	TextSearch string `json:"searchString" bson:"$text,omitempty"`
}

func (search *ImageSearch) Valid() bool {
	if !validMap(search.PublishTime) {
		return false
	} else if !validMap(search.Downloads) {
		return false
	} else if !validMap(search.MetaData.CaptureTime) {
		return false
	} else if !validMap(search.MetaData.Location) {
		return false
	}

	userRef, err := refs.GetRef(search.Owner)
	if err != nil {
		return false
	}
	search.OwnerExtern = userRef
	return true
}

func validMap(x interface{}) bool {
	m, ok := x.(map[Filter]interface{})
	if !ok {
		if len(m) == 0 {
			return true
		}
		return false
	}
	for key := range m {
		if !key.Valid() {
			return false
		}
	}
	return true
}
