package search

import (
	"fmt"
	"log"

	"github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/clr/clr"
	"github.com/fokal/fokal/pkg/handler"
	"github.com/lib/pq"
)

func color(state *handler.State, color clr.Color, pixelFraction float64, limit int) ([]Rank, error) {
	ids := []Rank{}

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
		return []Rank{}, err
	}

	return ids, nil
}

func geoRadius(state *handler.State, point postgis.PointS, radius float64, limit int) ([]Rank, error) {
	ids := []Rank{}

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
		return []Rank{}, err
	}

	log.Println(ids)
	return ids, nil
}
