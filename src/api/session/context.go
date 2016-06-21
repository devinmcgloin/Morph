package session

import (
	"net/http"

	"github.com/devinmcgloin/sprioc/src/model"
	"github.com/gorilla/context"
)

func GetUser(r *http.Request) (model.User, bool) {
	if user := context.Get(r, "authUser"); user != nil {
		return user.(model.User), true
	}
	return model.User{}, false
}

func SetUser(r *http.Request, user model.User) {
	context.Set(r, "authUser", user)
}
