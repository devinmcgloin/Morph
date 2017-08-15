package search

import (
	"log"
	"net/http"
	"strconv"

	"net/url"

	"strings"

	"sort"

	"github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/clr/clr"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/retrieval"
)

func limitParam(params url.Values) (limit int) {
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

func geoDistanceParams(params url.Values) bool {
	_, ok := params["lat"]
	if !ok {
		return false
	}

	_, ok = params["lng"]
	if !ok {
		return false
	}

	_, ok = params["radius"]
	if !ok {
		return false
	}

	return true

	//log.Printf("%+v %f %d\n", postgis.PointS{SRID: 4326, X: lng, Y: lat}, radius, limit)
	//images, err := GeoRadius(store, postgis.PointS{SRID: 4326, X: lng, Y: lat}, radius, limit)
	//if err != nil {
	//	return rsp, err
	//}
	//
	//return handler.Response{
	//	Code: http.StatusOK,
	//	Data: images,
	//}, nil
}

func textParams(params url.Values) bool {
	_, ok := params["q"]
	if !ok {
		return false
	}

	return true

	//q = strings.Split(q[0], " ")
	//
	//query := strings.Join(q, " | ")
	//
	//log.Printf("%d %+v %d %s\n", limit, q, len(q), query)
	//images, err := Text(store, query, limit)
	//if err != nil {
	//	return rsp, err
	//}
	//
	//return handler.Response{
	//	Code: http.StatusOK,
	//	Data: images,
	//}, nil
}

func colorParams(params url.Values) bool {

	hex, ok := params["hex"]
	if !ok {
		return false
	}

	if len(hex[0]) != 6 {
		return false
	}

	return true
}

func SearchHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	params := r.URL.Query()
	ranking := []Score{}
	limit := limitParam(params)

	if textParams(params) {
		q := params.Get("q")
		terms := strings.Split(q, " ")

		query := strings.Join(terms, " | ")

		r, err := text(store, query, limit)
		if err != nil {
			return rsp, err
		}

		ranking = append(ranking, r...)
	}

	if colorParams(params) {
		hex := params.Get("hex")
		pixelFraction := params.Get("pixel_fraction")
		pxl, err := strconv.ParseFloat(pixelFraction, 64)
		if pixelFraction == "" || err != nil {
			pxl = .005
		}

		r, err := color(store, clr.Hex{Code: hex}, pxl, limit)
		if err != nil {
			return rsp, err
		}
		ranking = append(ranking, r...)
	}

	if geoDistanceParams(params) {
		lat, err := strconv.ParseFloat(params.Get("lat"), 64)
		if err != nil {
			return rsp, err
		}
		lng, err := strconv.ParseFloat(params.Get("lng"), 64)
		if err != nil {
			return rsp, err
		}
		radius, err := strconv.ParseFloat(params.Get("radius"), 64)
		if err != nil {
			return rsp, err
		}

		p := postgis.PointS{SRID: 4326, X: lng, Y: lat}

		r, err := geoRadius(store, p, radius, limit)
		if err != nil {
			return rsp, err
		}
		ranking = append(ranking, r...)
	}

	m := make(map[int64]float64)
	for _, s := range ranking {
		score, ok := m[s.ID]
		if ok {
			m[s.ID] = score + s.Score
		} else {
			m[s.ID] = s.Score
		}
	}

	ranking = []Score{}

	for k, v := range m {
		ranking = append(ranking, Score{ID: k, Score: v})
	}

	sort.Sort(Scores(ranking))
	imgs, err := retrieval.GetImages(store, restrict(ranking, limit))
	if err != nil {
		return rsp, err
	}
	return handler.Response{
		Code: http.StatusOK,
		Data: imgs,
	}, nil

}

func restrict(s []Score, limit int) []int64 {
	var l []int64
	for i, score := range s {
		if i > limit {
			break
		} else {
			l = append(l, score.ID)
		}
	}
	return l
}

func RecentImageHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	params := r.URL.Query()
	limit := limitParam(params)

	log.Printf("%d\n", limit)
	images, err := RecentImages(store, limit)
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
	limit := limitParam(params)

	images, err := FeaturedImages(store, limit)
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
	limit := limitParam(params)

	images, err := Hot(store, limit)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}
