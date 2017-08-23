package random

import (
	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/model"
	"github.com/fokal/fokal/pkg/retrieval"
)

func Image(state *handler.State, u *int64) (model.Image, error) {
	var id int64
	var err error
	if u != nil {
		err = state.DB.Get(&id, "SELECT random_image($1);", *u)
	} else {
		err = state.DB.Get(&id, "SELECT random_image();")
	}

	if err != nil {
		return model.Image{}, err
	}

	return retrieval.GetImage(state, id)
}
