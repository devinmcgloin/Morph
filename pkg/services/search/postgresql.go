package search

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/devinmcgloin/clr/clr"
	"github.com/fokal/fokal-core/pkg/domain"
	"github.com/fokal/fokal-core/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type PGSearchService struct {
	db    *sqlx.DB
	user  domain.UserService
	tag   domain.TagService
	image domain.ImageService
}

func New(db *sqlx.DB, user domain.UserService, tag domain.TagService, image domain.ImageService) *PGSearchService {
	search := &PGSearchService{
		db:    db,
		user:  user,
		tag:   tag,
		image: image,
	}
	search.RefreshMaterializedView()
	return search
}

func (pgs *PGSearchService) RefreshMaterializedView() {
	tick := time.NewTicker(time.Minute * 10)
	go func() {
		for range tick.C {
			log.Println("Refreshing Materialized View")
			_, err := pgs.db.Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY searches;")
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func (pgs *PGSearchService) FullSearch(ctx context.Context, req domain.SearchRequest) (*[]domain.Rank, error) {

	if req.Color != nil && (len(req.Color.HexCode) != 7 || req.Color.HexCode[0] != '#') {
		log.Println(req.Color.HexCode)
		return nil, errors.New("invalid Hex Code")
	}

	var genericColor clr.Color

	if req.Color != nil {
		genericColor = clr.Hex{Code: req.Color.HexCode[1:7]}
		if !genericColor.Valid() {
			return nil, errors.New("clr invalid Hex Code")
		}
	}

	var ids []domain.Rank

	tsQuery := formatQueryString(req.RequiredTerms, req.OptionalTerms, req.ExcludedTerms)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := psql.Select("searches.searchable_id as ID", "searches.searchable_type as type").From("searches").
		LeftJoin("content.image_geo AS geo ON searches.searchable_id = geo.image_id").
		LeftJoin("content.image_color_bridge AS bridge ON searches.searchable_id = bridge.image_id").
		LeftJoin("content.colors AS colors ON bridge.color_id = colors.id").Where(sq.Eq{"searches.searchable_type": req.Types}).
		Options("DISTINCT ON (ID, type)")

	if tsQuery == "" {
		q = q.Column("0 AS rank")

	} else {
		q = q.Column("ts_rank_cd(term, to_tsquery(?), 32 /* rank/(rank+1) */) AS rank", tsQuery).Where("to_tsquery(?) @@ term", tsQuery)
	}

	if req.Geo != nil {
		geo := req.Geo

		q = q.Where(`ST_Covers(ST_MakeEnvelope(
        ?, ?,
        ?, ?, 
        ?), geo.loc) `, geo.SW.Lng, geo.SW.Lat, geo.NE.Lng, geo.NE.Lat, 4326)
	}

	if req.Color != nil {
		l, a, b := genericColor.CIELAB()
		c := fmt.Sprintf("(%f, %f, %f)", l, a, b)
		q = q.Where("? :: CUBE <-> colors.cielab < 50", c).Where(sq.Gt{"bridge.pixel_fraction": 0.05}).Column("? :: CUBE <-> colors.cielab as color_dist", c).GroupBy("colors.cielab")
	}

	q = q.GroupBy("searches.searchable_type", "searches.searchable_id", "searches.term").OrderBy("ID", "type")

	sqlString, args, err := q.ToSql()
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	err = pgs.db.Select(&ids, sqlString, args...)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	sort.Sort(ByRankColor(ids))

	return &ids, nil
}
