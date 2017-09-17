package main

import (
	"encoding/json"
	"log"
	"os"

	"flag"

	"context"

	"googlemaps.github.io/maps"
)

func main() {
	f := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(f)

	var lat, lng float64
	flag.Float64Var(&lat, "lat", 0, "Coordinate Latitude")
	flag.Float64Var(&lng, "lng", 0, "Coordinate Longitude")

	flag.Parse()
	if lat == 0 && lng == 0 {
		log.Fatalf("Either lat or lng were not provided.\n")
	}

	mapsClient, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_API_TOKEN")))
	if err != nil {
		log.Fatal(err)
	}

	geocodeRequest := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{Lat: lat, Lng: lng},
		//ResultType: []string{"point_of_interest", "natural_feature", "neighborhood", "premise", "airport"},
	}

	results, err := mapsClient.Geocode(context.Background(), geocodeRequest)
	if err != nil {
		log.Fatal(err)
	}

	d, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(d)
}
