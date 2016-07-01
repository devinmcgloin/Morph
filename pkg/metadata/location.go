package metadata

import (
	"log"
	"os"

	"gopkg.in/mgo.v2/bson"

	gj "github.com/sprioc/geojson"
	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/store"

	"golang.org/x/net/context"

	"googlemaps.github.io/maps"
)

var mapsClient *maps.Client

func init() {
	var err error
	mapsClient, err = maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_SECRET")))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

func SetLocation(point gj.Point) {

	if store.Exists("locations", bson.M{"bounds": bson.M{"$geoIntersects": bson.M{"$geometry": point}}}) {
		log.Println("Location already found")
		return
	} else if store.Exists("locations", bson.M{"bounds": bson.M{"$near": bson.M{"$geometry": point, "$maxDistance": 1000}}}) {
		log.Println("Location already found")
		return
	}

	log.Println("Finding new location")

	r := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: float64(point.Coordinates[1]),
			Lng: float64(point.Coordinates[0]),
		},
		LocationType: []maps.GeocodeAccuracy{maps.GeocodeAccuracyApproximate, maps.GeocodeAccuracyGeometricCenter},
		ResultType: []string{"point_of_interest", "airport", "natural_feature", "route",
			"neighborhood", "political"},
	}
	result, err := mapsClient.Geocode(context.Background(), r)
	if err != nil {
		log.Println(err)
		return
	}

	bounds := result[0].Geometry.Viewport

	location := model.Location{
		GoogleLoc: result[0],
		Bounds:    getBounds(bounds),
	}

	err = store.Create("locations", location)
	if err != nil {
		log.Println(err)
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
