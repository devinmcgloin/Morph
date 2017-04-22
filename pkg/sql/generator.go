package sql

import (
	"github.com/sprioc/composer/pkg/generator"
)

func generateSC(collection uint32) (string, error) {

	id := generator.RandString(12)

	var exist bool
	var err error

	for exist, err = ExistsImage(id); exist == true; exist, err = ExistsImage(id) {
		if err != nil {
			return "", err
		}
		id = generator.RandString(12)
	}
	return id, nil
}
