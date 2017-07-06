package generator

import (
	"github.com/devinmcgloin/fokal/pkg/retrieval"
	"github.com/jmoiron/sqlx"
)

func GenerateSC(db *sqlx.DB, collection uint32) (string, error) {

	id := RandString(12)

	var exist bool
	var err error

	for exist, err = retrieval.ExistsImage(db, id); exist == true; exist, err = retrieval.ExistsImage(db, id) {
		if err != nil {
			return "", err
		}
		id = RandString(12)
	}
	return id, nil
}
