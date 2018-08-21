package search

import (
	"errors"
	"net/http"

	"io/ioutil"

	"encoding/json"
	"log"

	"strings"

	"fmt"

	"sort"

	sq "github.com/Masterminds/squirrel"
	"github.com/devinmcgloin/clr/clr"
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/model"
	"github.com/fokal/fokal-core/pkg/retrieval"
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

	log.Println(searchReq)

	if searchReq.Color != nil && (len(searchReq.Color.HexCode) != 7 || searchReq.Color.HexCode[0] != '#') {
		log.Println(searchReq.Color.HexCode)
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("invalid Hex Code"),
			Code: http.StatusBadRequest}
	}

	var genericColor clr.Color

	if searchReq.Color != nil {
		genericColor = clr.Hex{Code: searchReq.Color.HexCode[1:7]}
		if !genericColor.Valid() {
			return handler.Response{}, handler.StatusError{
				Err:  errors.New("clr invalid Hex Code"),
				Code: http.StatusBadRequest}
		}
	}

	var ids []Rank

	tsQuery := formatQueryString(searchReq.RequiredTerms, searchReq.OptionalTerms, searchReq.ExcludedTerms)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := psql.Select("searches.searchable_id as ID", "searches.searchable_type as type").From("searches").
		LeftJoin("content.image_geo AS geo ON searches.searchable_id = geo.image_id").
		LeftJoin("content.image_color_bridge AS bridge ON searches.searchable_id = bridge.image_id").
		LeftJoin("content.colors AS colors ON bridge.color_id = colors.id").Where(sq.Eq{"searches.searchable_type": searchReq.Types}).
		Options("DISTINCT ON (ID, type)")

	if tsQuery == "" {
		q = q.Column("0 AS rank")

	} else {
		q = q.Column("ts_rank_cd(term, to_tsquery(?), 32 /* rank/(rank+1) */) AS rank", tsQuery).Where("to_tsquery(?) @@ term", tsQuery)
	}

	if searchReq.Geo != nil {
		geo := searchReq.Geo

		q = q.Where(`ST_Covers(ST_MakeEnvelope(
        ?, ?,
        ?, ?, 
        ?), geo.loc) `, geo.SW.Longitude, geo.SW.Longitude, geo.NE.Longitude, geo.NE.Latitude, 4326)
	}

	if searchReq.Color != nil {
		l, a, b := genericColor.CIELAB()
		c := fmt.Sprintf("(%f, %f, %f)", l, a, b)
		q = q.Where("? :: CUBE <-> colors.cielab < 50", c).Where(sq.Gt{"bridge.pixel_fraction": 0.05}).Column("? :: CUBE <-> colors.cielab as color_dist", c).GroupBy("colors.cielab")
	}

	q = q.GroupBy("searches.searchable_type", "searches.searchable_id", "searches.term").OrderBy("ID", "type")

	sqlString, args, err := q.ToSql()
	if err != nil {
		log.Println(err)
		return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
	}

	err = store.DB.Select(&ids, sqlString, args...)
	if err != nil {
		log.Println(err)
		return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
	}

	sort.Sort(ByRankColor(ids))

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
			tag, err := retrieval.TaggedImages(store, v.ID, 1)
			if err != nil {
				log.Println(err)
				return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
			}

			if len(tag.Images) == 0 {
				continue
			}
			ref, err := retrieval.GetTagRef(store.DB, v.ID)
			if err != nil {
				log.Println(err)
				return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
			}

			resp.Tags = append(resp.Tags, TagResponse{
				Id:         ref.Shortcode,
				Permalink:  ref.ToURL(store.Port, store.Local),
				TitleImage: tag.Images[0],
				Count:      tag.Count,
			})
		}
	}

	return handler.Response{Code: http.StatusOK, Data: resp}, nil

}

func formatQueryString(req []string, opt []string, exc []string) string {
	args := []string{}

	trim(req)
	req = FilterEmpty(req)
	trim(opt)
	opt = FilterEmpty(opt)
	trim(exc)
	exc = FilterEmpty(exc)

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

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func trim(vs []string) {
	for i, v := range vs {
		vs[i] = strings.Trim(v, ",./!@#$%^&*()_+-= ")
	}
}

func FilterEmpty(arr []string) []string {
	empty := func(s string) bool {
		return s != ""
	}

	return Filter(arr, empty)
}
