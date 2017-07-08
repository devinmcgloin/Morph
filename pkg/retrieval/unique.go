package retrieval

import (
	"github.com/devinmcgloin/fokal/pkg/generator"
	"github.com/jmoiron/sqlx"
)

func GenerateSC(db *sqlx.DB, collection uint32) (string, error) {

	id := generator.RandString(12)

	var exist bool
	var err error

	for exist, err = ExistsImage(db, id); exist == true; exist, err = ExistsImage(db, id) {
		if err != nil {
			return "", err
		}
		id = generator.RandString(12)
	}
	return id, nil
}
