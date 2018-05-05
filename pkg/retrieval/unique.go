package retrieval

import (
	"github.com/fokal/fokal-core/pkg/generator"
	"github.com/fokal/fokal-core/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func GenerateSC(db *sqlx.DB, collection model.ReferenceType) (string, error) {

	id := generator.RandString(12)

	var exist bool
	var err error

	var f func(db *sqlx.DB, id string) (bool, error)

	switch collection {
	case model.Images:
		f = ExistsImage
	default:
		return "", errors.New("Invalid Collection Type.")
	}

	for exist, err = f(db, id); exist; exist, err = f(db, id) {
		if err != nil {
			return "", err
		}
		id = generator.RandString(12)
	}
	return id, nil
}
