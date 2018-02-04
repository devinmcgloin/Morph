package routes

import (
	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/model"
	"github.com/fokal/fokal/pkg/security/permissions"
	"github.com/fokal/fokal/pkg/status"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterStatusRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {

	head := api.Methods("HEAD").Subrouter()
	head.Handle("/stauts",
		chain.Append(
			permissions.Middleware{State: state,
				T:          permissions.CanView,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle,
			}.Handler).Then(handler.Handler{State: state, H: status.StatusHandler}))

}
