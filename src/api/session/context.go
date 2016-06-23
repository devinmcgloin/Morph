package session

import (
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/context"
)

func GetUser(r *http.Request) (mgo.DBRef, bool) {
	if user := context.Get(r, "authUser"); user != nil {
		return user.(mgo.DBRef), true
	}
	return mgo.DBRef{}, false
}

func SetUser(r *http.Request, user mgo.DBRef) {
	context.Set(r, "authUser", user)
}
