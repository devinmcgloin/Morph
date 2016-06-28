package core

import (
	"errors"
	"strings"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/store"
)

func GetUser(ref model.DBRef) (model.User, error) {

	if strings.Compare(ref.Collection, "users") != 0 {
		return model.User{}, errors.New("Ref is of the wrong collection type")
	}

	doc, err := store.Get(ref)
	if err != nil {
		return model.User{}, errors.New("User not found")
	}

	user, ok := doc.(model.User)
	if !ok {
		return model.User{}, errors.New("Unable to cast document to user")
	}

	return user, nil
}

func GetImage(ref model.DBRef) (model.Image, error) {
	if strings.Compare(ref.Collection, "users") != 0 {
		return model.Image{}, errors.New("Ref is of the wrong collection type")
	}

	doc, err := store.Get(ref)
	if err != nil {
		return model.Image{}, errors.New("User not found")
	}

	user, ok := doc.(model.Image)
	if !ok {
		return model.Image{}, errors.New("Unable to cast document to user")
	}

	return user, nil
}

func GetCollection(ref model.DBRef) (model.Collection, error) {
	if strings.Compare(ref.Collection, "collections") != 0 {
		return model.Collection{}, errors.New("Ref is of the wrong collection type")
	}

	doc, err := store.Get(ref)
	if err != nil {
		return model.Collection{}, errors.New("User not found")
	}

	user, ok := doc.(model.Collection)
	if !ok {
		return model.Collection{}, errors.New("Unable to cast document to user")
	}

	return user, nil
}
