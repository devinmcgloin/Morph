package metadata

import (
	"context"
	"log"

	gj "github.com/sprioc/geojson"
	"googlemaps.github.io/maps"
)

func SetLocation(point *gj.Point) {
	r := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: float64(point.Coordinates[1]),
			Lng: float64(point.Coordinates[0]),
		},
		LocationType: []maps.GeocodeAccuracy{maps.GeocodeAccuracyApproximate, maps.GeocodeAccuracyGeometricCenter},
		ResultType: []string{"point_of_interest", "airport", "natural_feature", "route",
			"neighborhood", "political"},
	}
	_, err := mapsClient.Geocode(context.Background(), r)
	if err != nil {
		log.Println(err)
		return
	}
}

func getBounds(bounds maps.LatLngBounds) gj.Polygon {
	poly := &gj.Polygon{
		Type: "Polygon",
	}

	box := []gj.Coordinate{
		gj.Coordinate{gj.CoordType(bounds.NorthEast.Lng), gj.CoordType(bounds.NorthEast.Lat)},
		gj.Coordinate{gj.CoordType(bounds.SouthWest.Lng), gj.CoordType(bounds.NorthEast.Lat)},
		gj.Coordinate{gj.CoordType(bounds.SouthWest.Lng), gj.CoordType(bounds.SouthWest.Lat)},
		gj.Coordinate{gj.CoordType(bounds.NorthEast.Lng), gj.CoordType(bounds.SouthWest.Lat)},
		gj.Coordinate{gj.CoordType(bounds.NorthEast.Lng), gj.CoordType(bounds.NorthEast.Lat)},
	}

	poly.AddCoordinates(gj.Coordinates(box))
	return *poly
}
