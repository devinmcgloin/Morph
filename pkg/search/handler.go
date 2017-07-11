package search

import (
	"log"
	"net/http"
	"strconv"

	"net/url"

	"github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/clr/clr"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/gorilla/context"
)

func parsePaginationParams(params url.Values) (limit int, offset int, err error) {
	l, ok := params["limit"]

	if ok {
		if len(l) == 1 {
			limit, _ = strconv.Atoi(l[0])
		}
	}

	if limit == -1 {
		limit = 25
	}

	l, ok = params["offset"]
	if ok {
		if len(l) == 1 {
			offset, err = strconv.Atoi(l[0])
			if err != nil {
				offset = 0
			}
		}
	}

	if offset == -1 {
		offset = 0
	}

	return

}

func ColorHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	params := r.URL.Query()

	limit, offset, err := parsePaginationParams(params)
	if err != nil {
		return rsp, err
	}

	hex, ok := params["hex"]
	if !ok {
		return rsp, handler.StatusError{Code: http.StatusBadRequest}
	}

	var pixelFraction float64
	pixel, ok := params["pixelfraction"]
	if !ok {
		pixelFraction = .25
	} else {
		pixelFraction, err = strconv.ParseFloat(pixel[0], 64)
		if err != nil {
			pixelFraction = .25
		}
	}

	log.Printf("%d %d\n", limit, offset)
	images, err := Color(store, clr.Hex{Code: hex[0]}, pixelFraction, limit, offset)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}

func RecentImageHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	userRef, ok := context.GetOk(r, "auth")
	var usr int64
	if ok {
		usr = userRef.(model.Ref).Id
	}

	params := r.URL.Query()
	limit, offset, err := parsePaginationParams(params)
	if err != nil {
		return rsp, err

	}

	log.Printf("%d %d\n", limit, offset)
	images, err := RecentImages(store, usr, limit, offset)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}

func FeaturedImageHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	params := r.URL.Query()
	limit, offset, err := parsePaginationParams(params)
	if err != nil {
		return rsp, err

	}

	userRef, ok := context.GetOk(r, "auth")
	var usr int64
	if ok {
		usr = userRef.(model.Ref).Id
	}

	images, err := FeaturedImages(store, usr, limit, offset)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}

func HotImagesHander(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	params := r.URL.Query()
	limit, offset, err := parsePaginationParams(params)
	if err != nil {
		return rsp, err

	}

	images, err := Hot(store, limit, offset)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}

func GeoDistanceHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	params := r.URL.Query()

	limit, offset, err := parsePaginationParams(params)
	if err != nil {
		return rsp, err
	}

	var lat float64
	l, ok := params["lat"]
	if !ok {
		return rsp, handler.StatusError{Code: http.StatusBadRequest}
	} else {
		lat, err = strconv.ParseFloat(l[0], 64)
		if err != nil {
			return rsp, handler.StatusError{Code: http.StatusBadRequest}
		}
	}

	var lng float64
	l, ok = params["lng"]
	if !ok {
		return rsp, handler.StatusError{Code: http.StatusBadRequest}
	} else {
		lng, err = strconv.ParseFloat(l[0], 64)
		if err != nil {
			return rsp, handler.StatusError{Code: http.StatusBadRequest}
		}
	}

	var radius float64
	rad, ok := params["radius"]
	if !ok {
		radius = .25
	} else {
		radius, err = strconv.ParseFloat(rad[0], 64)
		if err != nil {
			radius = .25
		}
	}

	log.Printf("%d %d\n", limit, offset)
	images, err := GeoRadius(store, postgis.PointS{SRID: 4326, X: lat, Y: lng}, radius, limit, offset)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}
