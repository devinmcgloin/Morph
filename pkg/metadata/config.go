package metadata

import (
	"log"
	"net/http"

	"googlemaps.github.io/maps"

	"google.golang.org/api/googleapi/transport"
	vision "google.golang.org/api/vision/v1"
)

var visionService *vision.Service
var mapsClient *maps.Client

func Configure(apiKey string) {
	var err error

	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}
	visionService, err = vision.New(client)
	if err != nil {
		log.Fatal(err)
	}

	mapsClient, err = maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("fatal error: %s\n", err)
	}
}
