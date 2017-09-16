package geo

import (
	"context"
	"log"

	"errors"

	"github.com/cridenour/go-postgis"
	"github.com/fokal/fokal/pkg/handler"
	"googlemaps.github.io/maps"
)

func ReverseGeocode(state *handler.State, p *postgis.PointS) (string, error) {
	geocodeRequest := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{Lat: p.Y, Lng: p.X},
	}

	results, err := state.Maps.Geocode(context.Background(), geocodeRequest)
	if err != nil {
		log.Println(err)
		return "", err
	}

	for _, r := range results {
		return r.FormattedAddress, nil
	}
	return "", errors.New("no results found")
}
