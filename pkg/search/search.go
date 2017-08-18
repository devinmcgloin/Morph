package search

import (
	"fmt"

	"log"

	"github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/clr/clr"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/retrieval"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func color(state *handler.State, color clr.Color, pixelFraction float64, limit int) ([]Score, error) {
	ids := []Score{}

	l, a, b := color.CIELAB()
	cube := fmt.Sprintf("(%f, %f, %f)", l, a, b)

	log.Println(cube)
	err := state.DB.Select(&ids, `
	SELECT
	  id,
	  score / 50 AS score
	FROM (SELECT
			bridge.image_id                                         AS id,
			$1 :: CUBE <-> cielab AS score
		  FROM content.colors AS COLORS
			INNER JOIN content.image_color_bridge AS bridge ON COLORS.id = bridge.color_id
			INNER JOIN permissions.can_view AS view ON view.o_id = bridge.image_id
		  WHERE view.user_id = -1 AND bridge.pixel_fraction >= $3
		  ORDER BY score
		  LIMIT $2) AS scores
	WHERE score < 50
	GROUP BY id, score
	ORDER BY min(score);
	`, cube, limit, pixelFraction)
	if err != nil {
		return []Score{}, err
	}

	return ids, nil
}

func text(state *handler.State, query string, limit int) ([]Score, error) {
	ids := []Score{}

	err := state.DB.Select(&ids, `
SELECT
    scores.image_id AS id,
    sum(scores.rank) AS score
 FROM (SELECT
          bridge.image_id,
          ts_rank_cd(to_tsvector(landmark.description), to_tsquery($1),32 /* rank/(rank+1) */) AS rank
        FROM content.landmarks AS landmark
          JOIN content.image_landmark_bridge AS bridge ON landmark.id = bridge.landmark_id
        WHERE to_tsvector(landmark.description) @@ to_tsquery($1)
        UNION ALL
        SELECT
          bridge.image_id,
          ts_rank_cd(to_tsvector(labels.description), to_tsquery($1),32 /* rank/(rank+1) */) AS rank
        FROM content.labels AS labels
          JOIN content.image_label_bridge AS bridge ON labels.id = bridge.label_id
        WHERE to_tsvector(labels.description) @@ to_tsquery($1)
        UNION ALL
        SELECT
          bridge.image_id,
          ts_rank_cd(to_tsvector(tags.description), to_tsquery($1),32 /* rank/(rank+1) */) AS rank
        FROM content.image_tags AS tags
          JOIN content.image_tag_bridge AS bridge ON tags.id = bridge.tag_id
        WHERE to_tsvector(tags.description) @@ to_tsquery($1)) AS scores
  GROUP BY scores.image_id
  ORDER BY score DESC
  LIMIT $2;
	`, query, limit)
	if err != nil {
		return []Score{}, err
	}

	return ids, nil

}

func geoRadius(state *handler.State, point postgis.PointS, radius float64, limit int) ([]Score, error) {
	ids := []Score{}

	err := state.DB.Select(&ids, `
	SELECT geo.image_id AS id, ST_Distance(GeomFromEWKB($1), geo.loc) / $2 AS score
	FROM content.image_geo AS geo
	WHERE ST_Distance(GeomFromEWKB($1), geo.loc) < $2
	ORDER BY GeomFromEWKB($1) <-> geo.loc
	LIMIT $3
	`, point, radius, limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []Score{}, err
	}

	log.Println(ids)
	return ids, nil
}

func Hot(state *handler.State, limit int) ([]model.Image, error) {
	ids := []int64{}

	err := state.DB.Select(&ids, `
	SELECT id FROM content.images
	ORDER BY ranking(id, views + favorites , featured::int + 3) DESC
	LIMIT $1
	`, limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return retrieval.GetImages(state, ids)

}

func FeaturedImages(state *handler.State, limit int) ([]model.Image, error) {
	imgs := []int64{}
	var stmt *sqlx.Stmt
	var err error
	stmt, err = state.DB.Preparex(`
		SELECT images.id
		FROM content.images AS images
		INNER JOIN permissions.can_view AS view ON view.o_id = images.id
		WHERE view.user_id = -1 AND images.featured = TRUE
		ORDER BY publish_time DESC
		LIMIT $1
		`)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	err = stmt.Select(&imgs,
		limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return retrieval.GetImages(state, imgs)
}

func RecentImages(state *handler.State, limit int) ([]model.Image, error) {
	imageIds := []int64{}
	var err error
	err = state.DB.Select(&imageIds, `
		SELECT images.id
		FROM content.images AS images
		INNER JOIN permissions.can_view AS view ON view.o_id = images.id
		WHERE view.user_id = -1
		ORDER BY publish_time DESC
		LIMIT $1
		`, limit)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("GetRecentImages limit: %d %+v", limit, err)
		}
		return []model.Image{}, err
	}
	return retrieval.GetImages(state, imageIds)
}
