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

func Color(state *handler.State, color clr.Color, pixelFraction float64, limit, offset int) ([]model.Image, error) {
	ids := []int64{}

	l, a, b := color.CIELAB()
	cube := fmt.Sprintf("(%f, %f, %f)", l, a, b)

	err := state.DB.Select(&ids, `
			SELECT bridge.image_id
			FROM content.colors AS colors
  				INNER JOIN content.image_color_bridge AS bridge ON colors.id = bridge.color_id
  				INNER JOIN permissions.can_view AS view ON view.o_id = bridge.id
  				WHERE (view.user_id = -1 OR view.user_id = $1) AND bridge.pixel_fraction > $4
			ORDER BY $1::cube <-> cielab;
			OFFSET $2 LIMIT $3`, cube, offset, limit, pixelFraction)
	if err != nil {
		return []model.Image{}, err
	}

	return retrieval.GetImages(state, ids)
}

func Text(state *handler.State, query string, limit, offset int) ([]model.Image, error) {
	ids := []struct {
		Id   int64 `db:"image_id"`
		Rank float64
	}{}

	//images := []struct {
	//	Image model.Image
	//	Rank  float64
	//}{}

	err := state.DB.Select(&ids, `
			select text_search($1, $2, $3)`, query, offset, limit)
	if err != nil {
		return []model.Image{}, err
	}

	imageIds := make([]int64, len(ids))
	for i, v := range ids {
		imageIds[i] = v.Id
		//images[i].Rank = v.Rank
	}

	return retrieval.GetImages(state, imageIds)

}

func GeoRadius(state *handler.State, point postgis.PointS, radius float64, limit, offset int) ([]model.Image, error) {
	ids := []int64{}

	err := state.DB.Select(&ids, `
	SELECT geo.image_id
	FROM content.image_geo AS geo
	WHERE ST_Distance($1, geo.loc) < $2
	ORDER BY $1 <-> geo.loc
	OFFSET $3 LIMIT $4
	`, point, radius, offset, limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return retrieval.GetImages(state, ids)
}

func Hot(state *handler.State, limit, offset int) ([]model.Image, error) {
	ids := []int64{}

	err := state.DB.Select(&ids, `
	SELECT id FROM content.images
	ORDER BY ranking(id, views + favorites , featured::int + 3) DESC
	OFFSET $1 LIMIT $2
	`, offset, limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return retrieval.GetImages(state, ids)

}

func FeaturedImages(state *handler.State, userId int64, limit, offset int) ([]model.Image, error) {
	imgs := []int64{}
	var stmt *sqlx.Stmt
	var err error
	if userId == 0 {
		stmt, err = state.DB.Preparex(`
		SELECT images.id
		FROM content.images AS images
		INNER JOIN permissions.can_view AS view ON view.o_id = images.id
		WHERE (view.user_id = -1 OR view.user_id = $1) AND images.featured = TRUE
		ORDER BY publish_time DESC
		OFFSET $2 ROWS
		FETCH NEXT $3 ROWS ONLY
		`)
	} else {
		stmt, err = state.DB.Preparex(`
		SELECT images.id
		FROM content.images AS images
		INNER JOIN permissions.can_view AS view ON view.o_id = images.id
		WHERE view.user_id = -1 AND images.featured = TRUE
		ORDER BY publish_time DESC
		OFFSET $2 LIMIT $3
		`)
	}
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	err = stmt.Select(&imgs,
		userId, offset, limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return retrieval.GetImages(state, imgs)
}

func RecentImages(state *handler.State, userId int64, limit, offset int) ([]model.Image, error) {
	imageIds := []int64{}
	var stmt *sqlx.Stmt
	var err error
	if userId == 0 {
		stmt, err = state.DB.Preparex(`
		SELECT images.id
		FROM content.images AS images
		INNER JOIN permissions.can_view AS view ON view.o_id = images.id
		WHERE view.user_id = -1
		ORDER BY publish_time DESC
		OFFSET $2 LIMIT $3
		`)
	} else {
		stmt, err = state.DB.Preparex(`
		SELECT images.id
		FROM content.images AS images
		INNER JOIN permissions.can_view AS view ON view.o_id = images.id
		WHERE (view.user_id = -1 OR view.user_id = $1) AND images.user_id = $1
		ORDER BY publish_time DESC
		OFFSET $2 LIMIT $3
		`)
	}
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("GetRecentImages userId: %d limit: %d offset: %d %+v", userId, limit, offset, err)
		}
		return []model.Image{}, err
	}
	err = stmt.Select(&imageIds,
		userId, offset, limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("GetRecentImages userId: %d limit: %d offset: %d %+v", userId, limit, offset, err)
		}
		return []model.Image{}, err
	}
	log.Printf("GetRecentImages userId: %d limit: %d offset: %d %+v", userId, limit, offset, imageIds)

	return retrieval.GetImages(state, imageIds)
}
