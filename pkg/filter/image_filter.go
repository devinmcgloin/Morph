package filter

import (
	"fmt"
	"strings"
	"time"

	"github.com/sprioc/conductor/pkg/model"
	gj "github.com/sprioc/geojson"
)

var ordOps = []string{"$eq", "$gt", "$gte", "$lt", "$lte", "$ne"}

type ImageFilter struct {
	Where    string `bson:"$where" json:"-"`
	MetaData struct {
		Aperature       map[string]model.Ratio `bson:"-" json:"aperature,omitempty"`
		Exposure        map[string]model.Ratio `bson:"-" json:"exposure_time,omitempty"`
		FocalLength     map[string]model.Ratio `bson:"-" json:"focal_length,omitempty"`
		ISO             map[string]int         `bson:"iso,omitempty" json:"iso,omitempty"`
		Make            string                 `bson:"make,omitempty" json:"make,omitempty"`
		PixelXDimension map[string]int         `bson:"pixel_xd,omitempty" json:"pixel_xd,omitempty"`
		PixelYDimension map[string]int         `bson:"pixel_yd,omitempty" json:"pixel_yd,omitempty"`
		CaptureTime     map[string]time.Time   `bson:"capture_time,omitempty" json:"capture_time,omitempty"`
		Location        map[string]gj.Point    `bson:"location,omitempty" json:"location,omitempty"`
	} `bson:"metadata" json:"metadata"`
	PublishTime map[string]time.Time `bson:"publish_time,omitempty" json:"publish_time,omitempty"`

	CollectionLink string      `bson:"-" json:"collection,omitempty" `
	Collection     model.DBRef `bson:"collections,omitempty" json:"-"`

	FavoritedByLink string      `bson:"-" json:"favorited_by,omitempty"`
	FavoritedBy     model.DBRef `bson:"favorited_by,omitempty" json:"-"`

	Featured  bool           `json:"featured,omitempty"`
	Downloads map[string]int `json:"downloads,omitempty"`
}

func verifyOrdMap(m map[string]interface{}) bool {
	for _, key := range m {
		for _, op := range ordOps {
			if op == key {
				return true
			}
		}
	}
	return false
}

func verifyLocMap(m map[string]interface{}) bool {
	for _, key := range m {
		for _, op := range []string{"$near", "$geoWithin"} {
			if op == key {
				return true
			}
		}
	}
	return false
}

func prepareImageFilter(filter *ImageFilter) {
}

//js 1den * 2num [opt] 2den * 1num
func generateWhereClause(filter map[string]model.Ratio) string {
	var op string
	var clauses []string

	for opt, ratio := range filter {
		switch opt {
		case "$eq":
			op = "="
			break
		case "$gt":
			op = "<"
			break
		case "$gte":
			op = "<="
			break
		case "$lt":
			op = ">"
			break
		case "$lte":
			op = ">="
			break
		case "$ne":
			op = "!="
			break
		}

		// TODO QA this function, See notes.
		clauses = append(clauses, fmt.Sprintf("this.den * %d %s this.num * %d", ratio.Num, op, ratio.Den))
	}

	return strings.Join(clauses, "&&")
}
