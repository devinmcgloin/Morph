package geo

import (
	"context"
	"log"

	"errors"

	"github.com/cridenour/go-postgis"
	"googlemaps.github.io/maps"
)

func ReverseGeocode(m *maps.Client, p *postgis.PointS) (string, error) {
	geocodeRequest := &maps.GeocodingRequest{
		LatLng:     &maps.LatLng{Lat: p.Y, Lng: p.X},
		ResultType: []string{"point_of_interest", "natural_feature", "neighborhood"},
	}

	results, err := m.Geocode(context.Background(), geocodeRequest)
	if err != nil {
		log.Println(err)
		return "", err
	}

	for _, r := range results {
		return r.FormattedAddress, nil
	}
	return "", errors.New("no results found")
}
