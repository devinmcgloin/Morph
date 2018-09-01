package handler

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/fokal/fokal-core/pkg/log"
	"github.com/fokal/fokal-core/pkg/services/permission"

	"github.com/gorilla/mux"
)

// Permission represents the state needed for permission middleware to accept or deny requests
type Permission struct {
	*State
	T          permission.Scope
	TargetType permission.ResourceClass
	M          func(state *State, scope permission.Scope, TargetType permission.ResourceClass, next http.Handler) http.Handler
}

// Handler makes Permission http.Handler compliant by running the M method in Permission
func (m Permission) Handler(next http.Handler) http.Handler {
	return m.M(m.State, m.T, m.TargetType, next)
}

// PermissionMiddle implements the logic for deciding if a user can interact with a resource.
func PermissionMiddle(state *State, scope permission.Scope, resouceClass permission.ResourceClass, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sc := mux.Vars(r)["ID"]

		userID, ok := ctx.Value(log.UserIDKey).(uint64)
		if !ok {
			log.WithContext(ctx).Warn("unable to load value from context")
		}
		if scope != permission.CanView && !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var id uint64
		switch resouceClass {
		case permission.UserClass:
			log.WithContext(ctx).Debug("fetching user for permissions")
			if sc == "me" {
				id = userID
			} else {
				user, err := state.UserService.UserByUsername(ctx, sc)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				id = user.ID
			}
		case permission.StreamClass:
			log.WithContext(ctx).Debug("fetching stream for permissions")
			stream, err := state.StreamService.StreamByShortcode(ctx, sc)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			id = stream.ID
		case permission.ImageClass:
			log.WithContext(ctx).Debug("fetching image for permissions")
			img, err := state.ImageService.ImageByShortcode(ctx, sc)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			id = img.ID
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.WithContext(ctx).WithFields(logrus.Fields{
			"user-id":        userID,
			"scope":          scope,
			"resource-class": resouceClass,
			"resource-id":    id,
		}).Info("checking permissions for fields")

		valid, err := state.PermissionService.ValidScope(ctx, userID, id, resouceClass, scope)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !valid {
			log.WithContext(ctx).WithFields(logrus.Fields{
				"user-id":        userID,
				"scope":          scope,
				"resource-class": resouceClass,
				"resource-id":    id,
			}).Warn("user does not have access to resource")
		}

		if !valid && scope != permission.CanView {
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
