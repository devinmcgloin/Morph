package stats

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func AddStat(db *sqlx.DB, imageId int64, t string) error {
	_, err := db.Exec(`SELECT log_stat($1, $2)`, imageId, t)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
