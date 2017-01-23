package redis

import (
	"github.com/sprioc/composer/pkg/generator"
	"github.com/sprioc/composer/pkg/model"
)

func GenerateShortCode(colletion model.RString) (model.Ref, error) {

	id := generator.RandString(12)
	ref := model.Ref{Collection: colletion, ShortCode: id}

	var exist bool
	var err error

	for exist, err = Exists(ref); exist == true; exist, err = Exists(ref) {
		if err != nil {
			return ref, err
		}
		id = generator.RandString(12)
		ref.ShortCode = id
	}
	return ref, nil
}
