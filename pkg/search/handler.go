package search

import (
	"log"
	"net/http"
	"strconv"

	"net/url"

	"strings"

	"github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/clr/clr"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/gorilla/context"
	"github.com/pkg/errors"
)

func parsePaginationParams(params url.Values) (limit int) {
	var err error
	l, ok := params["limit"]
	if ok {
		if len(l) == 1 {
			limit, err = strconv.Atoi(l[0])
			if err != nil {
				limit = 500
			}
		}
	}

	if limit == 0 {
		limit = 500
	}

	return

}

func TextHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	params := r.URL.Query()

	limit := parsePaginationParams(params)

	q, ok := params["q"]
	if !ok {
		return rsp, handler.StatusError{Code: http.StatusBadRequest, Err: errors.New("Missing hex url parameter.")}
	}

	query := strings.Join(q, " | ")

	log.Printf("%d %+v %d %s\n", limit, q, len(q), query)
	images, err := Text(store, query, limit)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}

func ColorHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	var err error
	params := r.URL.Query()

	limit := parsePaginationParams(params)

	hex, ok := params["hex"]
	if !ok {
		return rsp, handler.StatusError{Code: http.StatusBadRequest, Err: errors.New("Missing hex url parameter.")}
	}

	if len(hex[0]) != 6 {
		return rsp, handler.StatusError{Code: http.StatusBadRequest, Err: errors.New("Invalid Hex code. Only codes following #000000 are accepted.")}

	}

	var pixelFraction float64
	pixel, ok := params["pixelfraction"]
	if !ok {
		pixelFraction = .005
	} else {
		pixelFraction, err = strconv.ParseFloat(pixel[0], 64)
		if err != nil {
			pixelFraction = .005
		}
	}

	log.Printf("%d %s\n", limit, hex[0])
	images, err := Color(store, clr.Hex{Code: hex[0]}, pixelFraction, limit)
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
	limit := parsePaginationParams(params)

	log.Printf("%d\n", limit)
	images, err := RecentImages(store, usr, limit)
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
	limit := parsePaginationParams(params)

	userRef, ok := context.GetOk(r, "auth")
	var usr int64
	if ok {
		usr = userRef.(model.Ref).Id
	}

	images, err := FeaturedImages(store, usr, limit)
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
	limit := parsePaginationParams(params)

	images, err := Hot(store, limit)
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
	var err error
	params := r.URL.Query()

	limit := parsePaginationParams(params)

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
		radius = 1000
	} else {
		radius, err = strconv.ParseFloat(rad[0], 64)
		if err != nil {
			radius = 1000
		}
	}

	log.Printf("%+v %f %d %d\n", postgis.PointS{SRID: 4326, X: lng, Y: lat}, radius, limit)
	images, err := GeoRadius(store, postgis.PointS{SRID: 4326, X: lng, Y: lat}, radius, limit)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}
