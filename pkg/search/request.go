package search

import (
	"net/http"

	"github.com/mholt/binding"
)

type queryRange struct {
	Value     interface{}
	Predicate predicate
}

type predicate string

const (
	GreaterThan   = predicate("gt")
	GreaterThanEq = predicate("gte")
	LessThan      = predicate("lt")
	LessThanEq    = predicate("lte")
	Eq            = predicate("eq")
)

type ImagesRequest struct {
	Query struct {
		Color *struct {
			H                  int
			S                  int
			V                  int
			HRange             int
			SRange             int
			VRange             int
			PixelFraction      float64
			PixelFractionRange float64
		}
		Terms *string
		Tags  *[]string
	}
	Filter struct {
		PublishTime  queryRange
		LastModified queryRange
		Featured     *bool
		Downloads    queryRange
		Views        queryRange
		Favorites    queryRange
		GeoNear      queryRange
		Username     *string
	}
}

func (cf *ImagesRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Query.Color.H:         "h",

	}
}
