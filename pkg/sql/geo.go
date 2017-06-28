package sql

import (
	"github.com/sprioc/composer/pkg/model"
	gj "github.com/sprioc/geojson"
)

func GeoNear(point gj.Point, distance uint16) ([]model.Ref, error) {
	return []model.Ref{}, nil
}
