package permissions

import (
	"net/http"

	"log"

	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/model"
	"github.com/fokal/fokal/pkg/retrieval"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type Middleware struct {
	*handler.State
	T          Permission
	TargetType model.ReferenceType
	M          func(state *handler.State, p Permission, TargetType model.ReferenceType, next http.Handler) http.Handler
}

func (m Middleware) Handler(next http.Handler) http.Handler {
	return m.M(m.State, m.T, m.TargetType, next)
}

func PermissionMiddle(state *handler.State, p Permission, TargetType model.ReferenceType, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["ID"]

		usr, ok := context.GetOk(r, "auth")

		if p != CanView {
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				log.Println("Auth params not set")
				return
			}

			if usr == nil {
				log.Println("User is nil")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}

		user, ok := usr.(model.Ref)
		if !ok {
			user = model.Ref{}
		} else {
			log.Println(user)
		}

		var tarRef model.Ref
		var err error

		switch TargetType {
		case model.Images:
			tarRef, err = retrieval.GetImageRef(state.DB, id)
		case model.Users:
			tarRef, err = retrieval.GetUserRef(state.DB, id)
		}

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		valid, err := Valid(state.DB, user.Id, p, tarRef.Id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !valid && p != CanView {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
