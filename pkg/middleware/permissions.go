package middleware

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/fokal/fokal-core/pkg/log"
	"github.com/fokal/fokal-core/pkg/services/permission"

	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/gorilla/mux"
)

// Permission represents the state needed for permission middleware to accept or deny requests
type Permission struct {
	*handler.State
	T          permission.Scope
	TargetType permission.ResourceClass
	M          func(state *handler.State, scope permission.Scope, TargetType permission.ResourceClass, next http.Handler) http.Handler
}

// Handler makes Permission http.Handler compliant by running the M method in Permission
func (m Permission) Handler(next http.Handler) http.Handler {
	return m.M(m.State, m.T, m.TargetType, next)
}

// PermissionMiddle implements the logrusic for deciding if a user can interact with a resource.
func PermissionMiddle(state *handler.State, scope permission.Scope, resouceClass permission.ResourceClass, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.ParseUint(mux.Vars(r)["ID"], 10, 64)
		if err != nil {
			logrus.Error(err)
		}
		userID, ok := ctx.Value(log.UserIDKey).(uint64)

		if scope != permission.CanView {
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		valid, err := state.PermissionService.ValidScope(ctx, userID, id, resouceClass, scope)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
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
