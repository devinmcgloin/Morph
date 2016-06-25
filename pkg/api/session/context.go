package session

import (
	"net/http"

	"github.com/devinmcgloin/sprioc/pkg/model"
	"github.com/gorilla/context"
)

func GetUser(r *http.Request) (model.DBRef, bool) {
	if user := context.Get(r, "authUser"); user != nil {
		return user.(model.DBRef), true
	}
	return model.DBRef{}, false
}

func SetUser(r *http.Request, user model.DBRef) {
	context.Set(r, "authUser", user)
}
