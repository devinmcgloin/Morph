package search

import (
	"errors"
	"net/http"

	"io/ioutil"

	"encoding/json"
	"log"

	"strings"

	"github.com/fokal/fokal/pkg/handler"
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

	log.Printf("%+v\n", searchReq)

	var ids []Rank

	q, args, err := sqlx.In(
		`
	SELECT searches.searchable_id as ID,
	 ts_rank_cd(term, query, 32 /* rank/(rank+1) */ ) AS rank,
	 searches.searchable_type as type
	 FROM searches, to_tsquery(?) query
	WHERE query @@ term AND searches.searchable_type IN (?)
	order by rank desc
	`, formatQueryString(searchReq.RequiredTerms, searchReq.OptionalTerms, searchReq.ExcludedTerms),
		searchReq.Types)
	if err != nil {
		log.Println(err)
		return handler.Response{}, handler.StatusError{Err: err, Code: http.StatusInternalServerError}
	}

	q = store.DB.Rebind(q)

	//log.Println(q, args)

	err = store.DB.Select(&ids, q, args...)
	return handler.Response{Code: http.StatusOK, Data: ids}, nil

}

func formatQueryString(req []string, opt []string, exc []string) string {
	args := make([]string, 3)
	for i, ex := range exc {
		exc[i] = "!" + ex
	}
	if len(req) != 0 {
		args[0] = "(" + strings.Join(req, " & ") + ")"
	}

	if len(opt) != 0 {
		args[1] = "(" + strings.Join(opt, " | ") + ")"
	}

	if len(exc) != 0 {
		args[2] = strings.Join(exc, " & ")
	}
	return strings.Join(args, " ")
}
