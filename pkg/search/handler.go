package search

import (
	"errors"
	"net/http"

	"io/ioutil"

	"encoding/json"
	"log"

	"strings"

	"fmt"

	"github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/clr/clr"
	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/model"
	"github.com/fokal/fokal/pkg/retrieval"
	"github.com/jmoiron/sqlx"
)

func SearchHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	file, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Unable to read request body."),
			Code: http.StatusBadRequest}
	}

	var searchReq Request
	err = json.Unmarshal(file, &searchReq)
	if err != nil {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Unable to parse request body."),
			Code: http.StatusBadRequest}
	}

	if searchReq.Color != nil && (len(searchReq.Color.HexCode) != 7 || searchReq.Color.HexCode[0] != '#') {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("invalid Hex Code"),
			Code: http.StatusBadRequest}
	}

	var ids []Rank
	var query string

	tsQuery := formatQueryString(searchReq.RequiredTerms, searchReq.OptionalTerms, searchReq.ExcludedTerms)
	var initialArgs []interface{}
	if tsQuery == "" {
		query = `
		SELECT
		  searches.searchable_id                          AS ID,
		  0 											  AS rank,
		  searches.searchable_type                        AS type
		FROM searches
		  LEFT JOIN content.image_geo AS geo ON searches.searchable_id = geo.image_id
		  LEFT JOIN content.image_color_bridge AS bridge ON searches.searchable_id = bridge.image_id
		  LEFT JOIN content.colors AS colors ON bridge.color_id = colors.id
		WHERE searches.searchable_type IN ( ? )`
		initialArgs = append(initialArgs, searchReq.Types)

	} else {
		query = `
		SELECT
		  searches.searchable_id                          AS ID,
		  ts_rank_cd(term, to_tsquery(?), 32 /* rank/(rank+1) */) AS rank,
		  searches.searchable_type                        AS type
		FROM searches
		  LEFT JOIN content.image_geo AS geo ON searches.searchable_id = geo.image_id
		  LEFT JOIN content.image_color_bridge AS bridge ON searches.searchable_id = bridge.image_id
		  LEFT JOIN content.colors AS colors ON bridge.color_id = colors.id
		WHERE to_tsquery(?) @@ term AND searches.searchable_type IN ( ? )`
		initialArgs = append(initialArgs, tsQuery, tsQuery, searchReq.Types)

	}

	log.Printf("TS_QUERY: {%s}\n", tsQuery)

	if searchReq.Geo != nil {
		query = query + `
		AND ST_Distance(GeomFromEWKB( ? ), geo.loc) < ?`
		geo := searchReq.Geo
		p := postgis.PointS{X: geo.Longitude, Y: geo.Latitude, SRID: 4326}
		initialArgs = append(initialArgs, p, geo.Radius)
	}

	if searchReq.Color != nil {
		query = query + `
		AND bridge.pixel_fraction >= ? AND ? :: CUBE <-> colors.cielab < 50`
		color := searchReq.Color
		genericColor := clr.Hex{Code: color.HexCode[1:7]}
		l, a, b := genericColor.CIELAB()
		c := fmt.Sprintf("(%f, %f, %f)", l, a, b)
		initialArgs = append(initialArgs, color.PixelFraction, c)
	}

	query = query + `
		GROUP BY searches.searchable_type, searches.searchable_id, searches.term
		ORDER BY rank DESC;`

	q, args, err := sqlx.In(query, initialArgs...)
	if err != nil {
		log.Println(err)
		return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
	}

	q = store.DB.Rebind(q)

	err = store.DB.Select(&ids, q, args...)
	if err != nil {
		log.Println(err)
		return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
	}

	resp := Response{
		Images: []model.Image{},
		Users:  []model.User{},
		Tags:   []TagResponse{}}

	for _, v := range ids {
		switch v.Type {
		case Image:
			img, err := retrieval.GetImage(store, v.ID)
			if err != nil {
				log.Println(err)
				return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
			}
			resp.Images = append(resp.Images, img)
		case User:
			user, err := retrieval.GetUser(store, v.ID)
			if err != nil {
				log.Println(err)
				return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
			}
			resp.Users = append(resp.Users, user)
		case Tag:
			img, err := retrieval.TaggedImages(store, v.ID, 1)
			if err != nil {
				log.Println(err)
				return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
			}

			if len(img) == 0 {
				continue
			}
			ref, err := retrieval.GetTagRef(store.DB, v.ID)
			if err != nil {
				log.Println(err)
				return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
			}
			resp.Tags = append(resp.Tags, TagResponse{Id: ref.Shortcode, Permalink: ref.ToURL(store.Port, store.Local), TitleImage: img[0]})
		}
	}

	return handler.Response{Code: http.StatusOK, Data: resp}, nil

}

func formatQueryString(req []string, opt []string, exc []string) string {
	args := []string{}
	for i, ex := range exc {
		exc[i] = "!" + ex
	}
	if len(req) != 0 {
		args = append(args, "("+strings.Join(req, " & ")+")")
	}

	if len(opt) != 0 {
		args = append(args, "("+strings.Join(opt, " | ")+")")
	}

	if len(exc) != 0 {
		args = append(args, strings.Join(exc, " & "))
	}
	return strings.Join(args, " & ")
}
