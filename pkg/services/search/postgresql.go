package search

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type PGSearchService struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PGSearchService {
	search := &PGSearchService{
		db: db,
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
